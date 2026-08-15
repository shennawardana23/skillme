# Laravel Patterns — Eloquent and API Reference

## Eager loading to avoid N+1

```php
$reservations = Reservation::query()
    ->with(['guest', 'hotel', 'roomAssignments.room'])
    ->latest()
    ->paginate(25);
```

## Explicit fillable (never $guarded = [])

```php
final class Reservation extends Model
{
    protected $fillable = ['guest_id', 'hotel_id', 'check_in', 'check_out', 'status'];

    protected $casts = [
        'status'    => ReservationStatus::class,
        'check_in'  => 'datetime',
        'check_out' => 'datetime',
    ];
}
```

## Scoped nested route bindings

```php
Route::middleware('auth:sanctum')->prefix('hotels/{hotel}')->group(function () {
    Route::scopeBindings()->group(function () {
        Route::get('/reservations/{reservation}', [ReservationController::class, 'show']);
    });
});
```

## Query scopes for reusable filters

```php
final class Reservation extends Model
{
    public function scopeForHotel(Builder $query, string $hotelId): Builder
    {
        return $query->where('hotel_id', $hotelId);
    }

    public function scopeActive(Builder $query): Builder
    {
        return $query->whereNull('cancelled_at');
    }
}

$reservations = Reservation::forHotel($hotelId)->active()->get();
```

## Form requests and DTOs

```php
final class StoreReservationRequest extends FormRequest
{
    public function authorize(): bool
    {
        return $this->user()?->can('create', Reservation::class) ?? false;
    }

    public function rules(): array
    {
        return [
            'guest_id'  => ['required', 'integer', 'exists:guests,id'],
            'check_in'  => ['required', 'date'],
            'check_out' => ['required', 'date', 'after:check_in'],
        ];
    }

    public function toDto(): CreateReservationData
    {
        return new CreateReservationData(
            guestId: (int) $this->validated('guest_id'),
            checkIn: $this->validated('check_in'),
            checkOut: $this->validated('check_out'),
        );
    }
}
```

## API resources with consistent pagination shape

```php
$reservations = Reservation::forHotel($hotelId)->active()->paginate(25);

return response()->json([
    'data' => ReservationResource::collection($reservations->items()),
    'meta' => [
        'page'     => $reservations->currentPage(),
        'per_page' => $reservations->perPage(),
        'total'    => $reservations->total(),
    ],
]);
```

## Idempotent queued jobs

```php
final class SendCheckInReminder implements ShouldQueue
{
    public function handle(): void
    {
        // idempotent: checking existence before acting survives redelivery
        if (Notification::where('reservation_id', $this->reservationId)
            ->where('type', 'check_in_reminder')
            ->exists()) {
            return;
        }
        // ... send notification, then record it
    }
}
```
