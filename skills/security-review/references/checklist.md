# Security Review — Extended Checklist

## Secrets management

```go
apiKey := os.Getenv("STRIPE_API_KEY")
if apiKey == "" {
	return fmt.Errorf("STRIPE_API_KEY not configured")
}
```

- [ ] No hardcoded keys/tokens/passwords anywhere in source
- [ ] `.env*` files in `.gitignore`; verify no secret ever landed in git history
- [ ] Production secrets live in the hosting platform's secret store, not in code

## Authorization (per-resource, not just per-endpoint)

```go
func DeleteReservation(ctx context.Context, requesterID, reservationID, hotelID string) error {
	reservation, err := repo.Get(ctx, hotelID, reservationID)
	if err != nil {
		return err
	}
	if !canManage(ctx, requesterID, reservation.HotelID) {
		return ErrForbidden // checked AFTER loading the specific resource, not just the role
	}
	return repo.Delete(ctx, hotelID, reservationID)
}
```

## Rate limiting

Apply stricter limits to expensive operations (search, export, password
reset) than to ordinary reads. Rate-limit by both IP and authenticated user
ID — IP-only limiting is trivially bypassed by rotating source addresses;
user-only limiting doesn't stop pre-auth abuse.

## CSRF and cookies

```
Set-Cookie: session=<id>; HttpOnly; Secure; SameSite=Strict
```

State-changing requests (POST/PUT/PATCH/DELETE) from a browser session
should carry a CSRF token validated server-side, independent of
`SameSite=Strict`, which mitigates but doesn't eliminate CSRF for all
browser/proxy configurations in the wild.

## Content Security Policy (browser-facing apps)

```
Content-Security-Policy: default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data: https:; frame-ancestors 'none';
```

Avoid `'unsafe-inline'`/`'unsafe-eval'` unless a specific, documented
dependency requires it — each weakens the policy for the entire origin.

## Dependency hygiene

```bash
# Go
go list -m -u all
govulncheck ./...

# Node
npm audit
npm ci   # not npm install, for reproducible CI builds

# PHP / Composer
composer audit
```

Commit lock files (`go.sum`, `package-lock.json`, `composer.lock`) — they
are what makes `npm audit`/`govulncheck` results reproducible across
machines and CI runs.
