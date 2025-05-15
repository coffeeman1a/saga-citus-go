package saga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/coffeeman1a/saga-citus-go/orchestrator/internal/config"
	"github.com/coffeeman1a/saga-citus-go/orchestrator/internal/models"
)

func HandleEvent(event *models.SagaEvent) error {
	// log start of step
	log.WithField("step", event.Step).
		WithFields(log.Fields{
			"user_id": event.UserID,
			"item":    event.Item,
			"price":   event.Price,
		}).
		Info("Handling saga step")

	switch event.Step {
	case "create_order":
		if err := callOrderService(event); err != nil {
			log.WithField("step", event.Step).
				WithError(err).
				Error("create_order failed")
			return finalizeOrder(event, "cancelled")
		}
		log.WithField("step", event.Step).
			WithField("order_id", event.OrderID).
			Info("Step success")
		return callReservePayment(event)

	case "reserve_payment":
		if err := callPaymentService(event); err != nil {
			log.WithField("step", event.Step).
				WithError(err).
				Error("reserve_payment failed")
			return finalizeOrder(event, "cancelled")
		}
		log.WithField("step", event.Step).
			WithField("order_id", event.OrderID).
			Info("Step success")
		return finalizeOrder(event, "accepted")

	default:
		err := fmt.Errorf("unknown step: %s", event.Step)
		log.WithField("step", event.Step).
			WithError(err).
			Error("invalid saga step")
		return err
	}
}

func callOrderService(event *models.SagaEvent) error {
	body := map[string]interface{}{
		"user_id": event.UserID,
		"item":    event.Item,
		"price":   event.Price,
	}
	log.WithFields(log.Fields{
		"step":         event.Step,
		"substep":      "callOrderService",
		"request_body": body,
	}).Info("Sending request to Order Service")

	data, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/orders", config.OrdersService)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.WithFields(log.Fields{
			"step":    event.Step,
			"substep": "callOrderService",
		}).
			WithError(err).
			Error("HTTP request failed to Order Service")
		return err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	log.WithFields(log.Fields{
		"step":        event.Step,
		"substep":     "callOrderService",
		"status_code": resp.StatusCode,
		"response":    string(respBytes),
	}).Info("Received response from Order Service")

	if resp.StatusCode >= 300 {
		err := fmt.Errorf("order-service returned status %s", resp.Status)
		log.WithFields(log.Fields{
			"step":        event.Step,
			"substep":     "callOrderService",
			"status_code": resp.StatusCode,
		}).
			WithError(err).
			Error("Order Service returned error status")
		return err
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		log.WithFields(log.Fields{
			"step":     event.Step,
			"substep":  "callOrderService",
			"response": string(respBytes),
		}).
			WithError(err).
			Error("Failed to parse Order Service response")
		return err
	}

	event.OrderID = result.ID
	log.WithFields(log.Fields{
		"step":     event.Step,
		"substep":  "callOrderService",
		"order_id": result.ID,
	}).Info("Order created")
	return nil
}

func callReservePayment(event *models.SagaEvent) error {
	event.Step = "reserve_payment"
	return HandleEvent(event)
}

func callPaymentService(event *models.SagaEvent) error {
	body := map[string]interface{}{
		"order_id": event.OrderID,
		"amount":   event.Price,
	}
	log.WithFields(log.Fields{
		"step":         event.Step,
		"substep":      "callPaymentService",
		"request_body": body,
	}).Info("Sending request to Payment Service")

	data, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/payments", config.PaymentsService)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.WithFields(log.Fields{
			"step":    event.Step,
			"substep": "callPaymentService",
		}).
			WithError(err).
			Error("HTTP request failed to Payment Service")
		return err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	log.WithFields(log.Fields{
		"step":        event.Step,
		"substep":     "callPaymentService",
		"status_code": resp.StatusCode,
		"response":    string(respBytes),
	}).Info("Received response from Payment Service")

	if resp.StatusCode >= 300 {
		err := fmt.Errorf("payment-service returned status %s", resp.Status)
		log.WithFields(log.Fields{
			"step":        event.Step,
			"substep":     "callPaymentService",
			"status_code": resp.StatusCode,
		}).
			WithError(err).
			Error("Payment Service returned error status")
		return err
	}

	log.WithFields(log.Fields{
		"step":    event.Step,
		"substep": "callPaymentService",
	}).Info("Payment reserved")
	return nil
}

func finalizeOrder(event *models.SagaEvent, status string) error {
	log.WithFields(log.Fields{
		"step":     "finalizeOrder",
		"order_id": event.OrderID,
		"status":   status,
	}).Info("Finalizing order")

	url := fmt.Sprintf("%s/orders/%s/status", config.OrdersService, event.OrderID)
	body := map[string]string{"status": status}
	log.WithFields(log.Fields{
		"step":         "finalizeOrder",
		"substep":      "finalizeOrderRequest",
		"url":          url,
		"request_body": body,
	}).Info("Sending finalizeOrder request")

	data, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		log.WithFields(log.Fields{
			"step":    "finalizeOrder",
			"substep": "finalizeOrderRequest",
		}).
			WithError(err).
			Error("Failed to create HTTP request for finalizeOrder")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"step":    "finalizeOrder",
			"substep": "finalizeOrderRequest",
		}).
			WithError(err).
			Error("HTTP request failed for finalizeOrder")
		return err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	log.WithFields(log.Fields{
		"step":        "finalizeOrder",
		"substep":     "finalizeOrderResponse",
		"status_code": resp.StatusCode,
		"response":    string(respBytes),
	}).Info("Received finalizeOrder response")

	if resp.StatusCode >= 300 {
		err := fmt.Errorf("finalizeOrder returned status %s", resp.Status)
		log.WithFields(log.Fields{
			"step":        "finalizeOrder",
			"substep":     "finalizeOrderResponse",
			"status_code": resp.StatusCode,
		}).
			WithError(err).
			Error("finalizeOrder error")
		return err
	}

	log.WithFields(log.Fields{
		"step":    "finalizeOrder",
		"substep": "finalizeOrder",
	}).Info("Order status updated successfully")
	return nil
}
