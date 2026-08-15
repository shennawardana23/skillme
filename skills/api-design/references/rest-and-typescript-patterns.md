# API Design — REST and TypeScript Patterns

## Resource design

```
GET    /api/reservations              list, paginated
POST   /api/reservations              create
GET    /api/reservations/:id          fetch one
PATCH  /api/reservations/:id          partial update
DELETE /api/reservations/:id          idempotent delete
```

## Pagination

```json
{
  "data": [ /* ... */ ],
  "pagination": { "page": 1, "pageSize": 20, "totalItems": 142, "totalPages": 8 }
}
```

## Discriminated unions for variant state (TypeScript)

```typescript
type ReservationStatus =
  | { type: "pending" }
  | { type: "confirmed"; confirmedAt: Date }
  | { type: "cancelled"; reason: string; cancelledAt: Date };

function label(status: ReservationStatus): string {
  switch (status.type) {
    case "pending": return "Pending";
    case "confirmed": return `Confirmed ${status.confirmedAt}`;
    case "cancelled": return `Cancelled: ${status.reason}`;
  }
}
```

## Branded IDs to prevent cross-entity mixups

```typescript
type HotelID = string & { readonly __brand: "HotelID" };
type ReservationID = string & { readonly __brand: "ReservationID" };

function getReservation(hotelID: HotelID, id: ReservationID): Promise<Reservation> { /* ... */ }
```

## Go handler boundary validation

```go
func CreateReservation(w http.ResponseWriter, r *http.Request) {
	var req CreateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error())
		return
	}
	// req is now trusted for the rest of this call chain
	reservation, err := service.Create(r.Context(), req)
	// ...
}
```
