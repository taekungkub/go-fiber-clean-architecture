package core

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const pagingLimitDefault = 10

func PagingRequest(c *fiber.Ctx, defaultLimit int64) (int64, int64) {

	// default
	page := int64(1)
	limit := defaultLimit

	if limit <= 0 {
		limit = pagingLimitDefault
	}

	// parse page
	if p, err := strconv.ParseInt(c.Query("page", "1"), 10, 64); err == nil && p > 0 {
		page = p
	}

	// parse limit
	if l, err := strconv.ParseInt(c.Query("limit", "10"), 10, 64); err == nil && l > 0 {
		limit = l
	}

	return page, limit
}

type PagingParams struct {
	Page  int64
	Limit int64
	Sort  string
	Order string
}

type pagingBody struct {
	Page  any    `json:"page"`
	Limit any    `json:"limit"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

func PagingRequestBody(c *fiber.Ctx) PagingParams {
	params := PagingParams{
		Page:  1,
		Limit: pagingLimitDefault,
		Sort:  "",
		Order: "",
	}

	var body pagingBody
	if err := c.BodyParser(&body); err != nil {
		return params
	}

	params.Page = toInt64(body.Page, 1)
	params.Limit = toInt64(body.Limit, pagingLimitDefault)

	if body.Sort != "" {
		params.Sort = body.Sort
	}
	if body.Order != "" {
		params.Order = body.Order
	}

	return params
}

func toInt64(v any, def int64) int64 {
	switch x := v.(type) {
	case float64:
		if x > 0 {
			return int64(x)
		}
	case string:
		if i, err := strconv.ParseInt(x, 10, 64); err == nil && i > 0 {
			return i
		}
	}
	return def
}
