# Integration Tests

This directory contains both unit tests and integration tests for the services layer.

## Running Tests

### Unit Tests Only (default)
```bash
go test ./internal/services
```

### Integration Tests Only
```bash
go test ./internal/services -tags=integration
```

### All Tests (unit + integration)
```bash
go test ./internal/services -tags=integration
```

### With Coverage
```bash
# Unit tests only
go test ./internal/services -coverprofile=coverage.out
go tool cover -html=coverage.out

# Integration tests
go test ./internal/services -tags=integration -coverprofile=coverage-integration.out
go tool cover -html=coverage-integration.out
```

## Prerequisites

### R2 Integration Tests
The R2 integration tests require:

1. **Environment variables** set in your `.env` file:
   ```
   R2_ACCESS_KEY_ID=your-access-key
   R2_SECRET_ACCESS_KEY=your-secret-key
   R2_ACCOUNT_ID=your-account-id
   R2_BUCKET_NAME=your-bucket-name
   R2_PUBLIC_HOST=https://your-cdn.example.com  # optional
   ```

2. **Valid Cloudflare R2 credentials** with permissions to:
   - Upload files (PutObject)
   - Delete files (DeleteObject, DeleteObjects)
   - Generate presigned URLs

3. **Test cleanup**: Integration tests automatically clean up test files, but they create files under the `integration-tests/` folder in your bucket.

**What the tests do:**
- Upload test files (small text files)
- Delete files individually and in batches
- Generate presigned upload URLs
- Test metadata attachment

**Note:** Tests will skip automatically if R2 credentials are not configured.

### FFmpeg Integration Tests
The ffmpeg integration tests require:

1. **FFmpeg installed** on your system:
   ```bash
   # macOS
   brew install ffmpeg

   # Ubuntu/Debian
   sudo apt-get install ffmpeg

   # Check installation
   ffmpeg -version
   ```

**What the tests do:**
- Create minimal valid WAV audio files
- Stitch multiple WAV files into a single MP3
- Verify output file creation and size
- Test error handling (missing files, empty input)

**Note:** Tests will skip automatically if ffmpeg is not found.

## Integration Test Structure

Integration tests are marked with the build tag:
```go
//go:build integration
// +build integration
```

This allows you to:
- Run unit tests quickly in development (no external dependencies)
- Run integration tests in CI/CD or when testing with real services
- Keep integration test code separate from unit test code

## CI/CD Considerations

For CI/CD pipelines, you can:

1. **Run unit tests on every PR** (fast, no credentials needed):
   ```bash
   go test ./internal/services
   ```

2. **Run integration tests on main branch** (requires secrets):
   ```bash
   go test ./internal/services -tags=integration
   ```

3. **Use separate test buckets** for different environments:
   - Development: `briefbot-dev`
   - Staging: `briefbot-staging`
   - CI: `briefbot-ci-tests`

## Test Files Cleanup

R2 integration tests create files in the `integration-tests/` folder. If tests fail and don't clean up properly, you can manually delete them:

```bash
# List test files (if you have AWS CLI configured for R2)
aws s3 ls s3://your-bucket/integration-tests/ --endpoint-url=https://your-account-id.r2.cloudflarestorage.com

# Delete test files
aws s3 rm s3://your-bucket/integration-tests/ --recursive --endpoint-url=https://your-account-id.r2.cloudflarestorage.com
```

## Troubleshooting

### R2 Tests Failing
- Verify credentials are correct in `.env`
- Check that your R2 bucket exists
- Ensure your API token has sufficient permissions
- Check network connectivity to Cloudflare

### FFmpeg Tests Failing
- Verify ffmpeg is installed: `which ffmpeg`
- Check ffmpeg version: `ffmpeg -version`
- Ensure ffmpeg has necessary codecs (libmp3lame for MP3)

### Tests Timing Out
- Increase test timeout: `go test -timeout 5m -tags=integration`
- Check network speed for R2 uploads
- Verify R2 endpoint is accessible
