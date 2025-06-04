package controllers

import (
	e "errors"
	"net/http"
	"strconv"

	"github.com/hutamy/invoice-generator/services"
	"github.com/hutamy/invoice-generator/utils"
	"github.com/hutamy/invoice-generator/utils/errors"
	"github.com/labstack/echo/v4"
)

type PDFController struct {
	pdfService services.PDFService
}

func NewPDFController(pdfService services.PDFService) *PDFController {
	return &PDFController{pdfService: pdfService}
}

func (c *PDFController) DownloadInvoicePDF(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return utils.Response(ctx, http.StatusBadRequest, errors.ErrBadRequest.Error(), nil)
	}

	pdfData, err := c.pdfService.GenerateInvoicePDF(uint(id))
	if err != nil {
		if e.Is(err, errors.ErrNotFound) {
			return utils.Response(ctx, http.StatusNotFound, err.Error(), nil)
		}

		return utils.Response(ctx, http.StatusInternalServerError, "Failed to generate PDF", nil)
	}

	return ctx.Blob(http.StatusOK, "application/pdf", pdfData)
}
