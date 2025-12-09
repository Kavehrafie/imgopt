# Image Optimization Microservice

This service optimizes and resizes images stored in Backblaze B2.

## Requirements

- Go 1.23+
- [libvips](https://github.com/libvips/libvips) 8.10+ (Required for `bimg`)
- Docker (optional, for containerized deployment)

## Configuration

Create a `.env` file based on `.env.example`:

```dotenv
PORT=8080
B2_BUCKET_NAME=your-bucket-name
B2_ACCOUNT_ID=your-account-id
B2_APPLICATION_KEY=your-application-key
```

## Running Locally

**Note:** You must have `libvips` installed and `pkg-config` configured to run locally because `bimg` uses CGO.

```bash
go run ./cmd/server
```

## Running with Docker

The Dockerfile is set up to install all necessary dependencies.

```bash
docker build -t imgopt .
docker run -p 8080:8080 --env-file .env imgopt
```

## API

### Resize Image

```
GET /path/to/image.jpg/w_300/h_200/fit_cover/crop_center
```

- **Path**: The first part of the URL is the path to the image in the B2 bucket.
- **Options**: Append options as path segments after the image key.
  - `w_<number>`: Width (e.g., `w_300`)
  - `h_<number>`: Height (e.g., `h_200`)
  - `fit_<mode>`: Resize mode (e.g., `fit_cover`)
    - `cover` (default)
    - `contain`
    - `fill`
    - `inside`
  - `crop_<mode>`: Crop gravity (e.g., `crop_smart`)
    - `center` (default)
    - `top`, `bottom`, `left`, `right`, `smart`

**Examples:**

- Original image: `/my-folder/image.jpg`
- Resize to 300x200: `/my-folder/image.jpg/w_300/h_200`
- Smart crop to 500x500: `/image.jpg/w_500/h_500/fit_cover/crop_smart`

If `w` and `h` are omitted, the original image is returned (proxied).
