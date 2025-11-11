package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/fresh132/URL-check/internal"
	"github.com/fresh132/URL-check/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func Check(c *gin.Context) {
	var req struct {
		Url []string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || len(req.Url) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var wg sync.WaitGroup

	statuses := make([]internal.URLstatus, len(req.Url))

	for i, url := range req.Url {
		wg.Add(1)

		go func(idx int, u string) {
			defer wg.Done()

			statuses[idx] = internal.URLstatus{URL: u, Status: config.CheckURL(u)}
		}(i, url)
	}

	wg.Wait()

	internal.Mutx.Lock()

	id := internal.ID

	internal.ID++

	internal.Data[id] = internal.TimeURL{Link: statuses, Time: time.Now()}

	internal.Mutx.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"statuses": statuses,
		"id":       id,
	})

	config.Save()
}

func Report(c *gin.Context) {
	var req struct {
		Links_num []string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	internal.Mutx.Lock()
	var reports []internal.TimeURL
	var ids []int

	if len(req.Links_num) == 0 || (len(req.Links_num) == 1 && req.Links_num[0] == "") {
		for id, batch := range internal.Data {
			reports = append(reports, batch)
			ids = append(ids, id)
		}
	} else {
		for _, strid := range req.Links_num {
			intid, err := strconv.Atoi(strid)

			if err != nil {
				internal.Mutx.Unlock()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
				return
			}
			if batch, ok := internal.Data[intid]; ok {
				reports = append(reports, batch)
				ids = append(ids, intid)
			}
		}
	}
	internal.Mutx.Unlock()

	if len(reports) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=report.pdf")

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(0, 0, 0)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "URL Check Report")
	pdf.Ln(20)

	for i, batch := range reports {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, fmt.Sprintf("ID: %d, time: %s", ids[i], batch.Time.Format(time.RFC3339)))
		pdf.Ln(10)
		pdf.SetFont("Arial", "", 12)
		for _, ls := range batch.Link {
			pdf.Cell(40, 10, fmt.Sprintf("%s: %s", ls.URL, ls.Status))
			pdf.Ln(10)
		}
		pdf.Ln(10)
	}

	if err := pdf.Output(c.Writer); err != nil {
		log.Println("PDF generation error:", err)
	}
}
