package handlers

import (
	"context"
	"net/http"

	"github.com/cyberbono3/infura-coding/client-api/data"
)

func (c *ClientHandler) middlewareValidation(d, key interface{}, rw http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	//deserialize the data
	err := data.FromJSON(d, r.Body)
	if err != nil {
		c.l.Error("Deserializing data", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return nil, nil
	}

	// validate data

	errs := c.v.Validate(d)
	if len(errs) != 0 {
		c.l.Error("Validating data", "error", errs)

		// return the validation messages as an array
		rw.WriteHeader(http.StatusUnprocessableEntity)
		data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
		return nil, nil
	}

	ctx := context.WithValue(r.Context(), key, d)
	r = r.WithContext(ctx)

	return rw, r

}

func (c *ClientHandler) MiddlewareValidateTransactionData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		rw, r = c.middlewareValidation(&data.TransactionData{}, KeyTransaction{}, rw, r)
		if rw == nil || r == nil {
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)

	})
}

func (c *ClientHandler) MiddlewareValidateBlockData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		rw, r = c.middlewareValidation(&data.BlockData{}, KeyBlock{}, rw, r)
		if rw == nil || r == nil {
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
