package controllers

import (
	"net/http"
	"strconv"
	"time"

	e "errors"

	"github.com/hutamy/invoice-generator-backend/dto"
	"github.com/hutamy/invoice-generator-backend/models"
	"github.com/hutamy/invoice-generator-backend/services"
	"github.com/hutamy/invoice-generator-backend/utils"
	"github.com/hutamy/invoice-generator-backend/utils/errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type InvoiceController struct {
	invoiceService services.InvoiceService
}

func NewInvoiceController(invoiceService services.InvoiceService) *InvoiceController {
	return &InvoiceController{invoiceService: invoiceService}
}

// @Summary      Create a new invoice
// @Description  Creates a new invoice for the authenticated user
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        invoice  body      dto.CreateInvoiceRequest  true  "Invoice data"
// @Success      201      {object}  utils.GenericResponse
// @Failure      400      {object}  utils.GenericResponse
// @Failure      500      {object}  utils.GenericResponse
// @Router       /v1/protected/invoices [post]
func (c *InvoiceController) CreateInvoice(ctx echo.Context) error {
	var req dto.CreateInvoiceRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	if err := ctx.Validate(req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	userID := ctx.Get("user_id").(uint)
	dueDate, err := time.Parse(time.DateOnly, req.DueDate)
	if err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrInvalidDateFormat.Error(), nil)
	}

	issueDate, err := time.Parse(time.DateOnly, req.IssueDate)
	if err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrInvalidDateFormat.Error(), nil)
	}

	invoice := models.Invoice{
		UserID:        userID,
		InvoiceNumber: req.InvoiceNumber,
		ClientID:      req.ClientID,
		IssueDate:     issueDate,
		DueDate:       dueDate,
		Notes:         req.Notes,
		Status:        "draft", // default status
		TaxRate:       req.TaxRate,
		ClientName:    req.ClientName,
		ClientEmail:   req.ClientEmail,
		ClientAddress: req.ClientAddress,
		ClientPhone:   req.ClientPhone,
	}

	for _, item := range req.Items {
		invoice.Items = append(invoice.Items, models.InvoiceItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		})
	}
	if err := c.invoiceService.CreateInvoice(&invoice); err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.Response(ctx, http.StatusCreated, "Invoice created successfully", invoice)
}

// @Summary      Get invoice by ID
// @Description  Retrieves an invoice by its ID
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Invoice ID"
// @Success      200  {object}  utils.GenericResponse
// @Failure      400  {object}  utils.GenericResponse
// @Failure      404  {object}  utils.GenericResponse
// @Failure      500  {object}  utils.GenericResponse
// @Router       /v1/protected/invoices/{id} [get]
func (c *InvoiceController) GetInvoiceByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	invoice, err := c.invoiceService.GetInvoiceByID(uint(id))
	if err != nil {
		if err == errors.ErrNotFound {
			return utils.Response(ctx, http.StatusNotFound, err.Error(), nil)
		}
		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "Invoice retrieved successfully", invoice)
}

// @Summary      List invoices by user
// @Description  Retrieves all invoices for the authenticated user
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.GenericResponse
// @Failure      500  {object}  utils.GenericResponse
// @Router       /v1/protected/invoices [get]
func (c *InvoiceController) ListInvoicesByUserID(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uint)

	invoices, err := c.invoiceService.ListInvoiceByUserID(userID)
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "Invoices retrieved successfully", invoices)
}

// @Summary      Update an invoice
// @Description  Updates an invoice by its ID
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int                      true  "Invoice ID"
// @Param        invoice body      dto.UpdateInvoiceRequest  true  "Invoice data"
// @Success      200     {object}  utils.GenericResponse
// @Failure      400     {object}  utils.GenericResponse
// @Failure      404     {object}  utils.GenericResponse
// @Failure      500     {object}  utils.GenericResponse
// @Router       /v1/protected/invoices/{id} [put]
func (c *InvoiceController) UpdateInvoice(ctx echo.Context) error {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid invoice ID"})
	}

	var req dto.UpdateInvoiceRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}
	if err := ctx.Validate(&req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	if err := c.invoiceService.UpdateInvoice(uint(id), &req); err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return utils.Response(ctx, http.StatusNotFound, err.Error(), nil)
		}

		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "Invoice updated successfully", nil)
}

// DownloadInvoicePDF godoc
// @Summary      Download invoice PDF
// @Description  Generates and downloads the PDF for a given invoice ID
// @Tags         invoices
// @Produce      application/pdf
// @Param        id   path      int  true  "Invoice ID"
// @Success      200  {file}    file
// @Failure      400  {object}  utils.GenericResponse
// @Failure      404  {object}  utils.GenericResponse
// @Failure      500  {object}  utils.GenericResponse
// @Router       /v1/protectted/invoices/{id}/pdf [post]
func (c *InvoiceController) DownloadInvoicePDF(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	pdfData, err := c.invoiceService.GenerateInvoicePDF(uint(id))
	if err != nil {
		if e.Is(err, errors.ErrNotFound) {
			return utils.Response(ctx, http.StatusNotFound, err.Error(), nil)
		}

		return utils.Response(ctx, http.StatusInternalServerError, "Failed to generate PDF", nil)
	}

	return ctx.Blob(http.StatusOK, "application/pdf", pdfData)
}

// DownloadInvoicePDF godoc
// @Summary      Download invoice PDF
// @Description  Generates and downloads the PDF for a given invoice ID
// @Tags         invoices
// @Produce      application/pdf
// @Param        invoice body   dto.GeneratePublicInvoiceRequest  true  "Invoice data"
// @Success      200  {file}    file
// @Failure      400  {object}  utils.GenericResponse
// @Failure      404  {object}  utils.GenericResponse
// @Failure      500  {object}  utils.GenericResponse
// @Router       /v1/public/invoices/generate-pdf [post]
func (c *InvoiceController) GeneratePublicInvoice(ctx echo.Context) error {
	var req dto.GeneratePublicInvoiceRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	if err := ctx.Validate(req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	pdfData, err := c.invoiceService.GeneratePublicInvoicePDF(req)
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, "Failed to generate PDF", nil)
	}

	return ctx.Blob(http.StatusOK, "application/pdf", pdfData)
}

// @Summary      Delete an invoice
// @Description  Deletes an invoice by its ID
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Invoice ID"
// @Success      200  {object}  utils.GenericResponse
// @Failure      400  {object}  utils.GenericResponse
// @Failure      404  {object}  utils.GenericResponse
// @Failure      500  {object}  utils.GenericResponse
// @Router       /v1/protected/invoices/{id} [delete]
func (c *InvoiceController) DeleteInvoice(ctx echo.Context) error {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid invoice ID"})
	}

	if err := c.invoiceService.DeleteInvoice(uint(id)); err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return utils.Response(ctx, http.StatusNotFound, err.Error(), nil)
		}

		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "Invoice deleted successfully", nil)
}

// @Summary      Update invoice status
// @Description  Updates the status of an invoice by its ID
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int     true  "Invoice ID"
// @Param        status body      string  true  "New status for the invoice"
// @Success      200    {object}  utils.GenericResponse
// @Failure      400    {object}  utils.GenericResponse
// @Failure      404    {object}  utils.GenericResponse
// @Failure      500    {object}  utils.GenericResponse
// @Router       /v1/protected/invoices/{id}/status [put]

func (c *InvoiceController) UpdateInvoiceStatus(ctx echo.Context) error {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid invoice ID"})
	}

	var req dto.UpdateInvoiceStatusRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	if err := c.invoiceService.UpdateInvoiceStatus(uint(id), req.Status); err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			return utils.Response(ctx, http.StatusNotFound, err.Error(), nil)
		}

		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "Invoice status updated successfully", nil)
}

// @Summary      Get invoice summary
// @Description  Retrieves a summary of invoices for the authenticated user
// @Tags         invoices
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.GenericResponse
// @Failure      500  {object}  utils.GenericResponse
// @Router       /v1/protected/invoices/summary [get]
func (c *InvoiceController) InvoiceSummary(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uint)
	summary, err := c.invoiceService.InvoiceSummary(userID)
	if err != nil {
		return utils.Response(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.Response(ctx, http.StatusOK, "Invoice summary retrieved successfully", summary)
}
