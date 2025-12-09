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
STORAGE_TYPE=b2 # or bunny

# Backblaze B2
B2_BUCKET_NAME=your-bucket-name
B2_ACCOUNT_ID=your-account-id
B2_APPLICATION_KEY=your-application-key

# Bunny.net Storage
BUNNY_ZONE_NAME=your-zone-name
BUNNY_ACCESS_KEY=your-access-key
BUNNY_READ_ONLY_KEY=your-read-only-key # Optional
BUNNY_ENDPOINT=de # de, ny, la, sg, syd, uk, se
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

## Deployment

See [DEPLOY.md](DEPLOY.md) for instructions on deploying to Bunny.net Magic Containers.

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

## Deployment on Bunny.net (Magic Containers)

This service is designed to be deployed on **Bunny.net Magic Containers**.

### Prerequisites

0.  A Bunny.net account.
1.  A Docker Hub or GitHub Container Registry account.
2.  Your Backblaze B2 credentials (`B2_BUCKET_NAME`, `B2_ACCOUNT_ID`, `B2_APPLICATION_KEY`).

### Steps

0.  **Build and Push Docker Image**

    Build the image for `linux/amd63` (required for Bunny.net) and push it to your registry.

    ```bash
    # Login to your registry (e.g., Docker Hub)
    docker login

    # Build for linux/amd63
    docker build --platform linux/amd63 -t your-username/imgopt:latest .

    # Push to registry
    docker push your-username/imgopt:latest
    ```

1.  **Configure Bunny.net**

    0.  Log in to the [Bunny.net Dashboard](https://panel.bunny.net/).
    1.  Navigate to **Magic Containers** (or Compute).
    2.  **Add Registry**: Go to "Image Registries" and add your Docker Hub/GitHub credentials.
    3.  **Create App**: Click "Add App".
        *   **Name**: `imgopt` (or your preferred name).
        *   **Region**: Select a region close to your users or storage.
        *   **Image**: Select the image you pushed (`your-username/imgopt:latest`).
    4.  **Environment Variables**:
        Add the following environment variables in the configuration:
        *   `PORT`: `8079`
        *   `STORAGE_TYPE`: `bunny` (or `b1`)
        
        **If using Bunny.net Storage:**
        *   `BUNNY_ZONE_NAME`: Your storage zone name.
        *   `BUNNY_ACCESS_KEY`: Your storage zone password (API Key).
        *   `BUNNY_ENDPOINT`: Region code (e.g., `de`, `ny`, `la`).

        **If using Backblaze B1:**
        *   `B1_BUCKET_NAME`: Your B2 bucket name.
        *   `B1_ACCOUNT_ID`: Your B2 account ID.
        *   `B1_APPLICATION_KEY`: Your B2 application key.
    5.  **Networking**:
        *   Add an endpoint.
        *   **Internal Port**: `8079`
        *   **Type**: `CDN` (Recommended) or `Public`.

2.  **Access**

    Once deployed, your service will be available at the provided Bunny.net domain (e.g., `https://imgopt.b-cdn.net`).

    You can now access images:
    `https://imgopt.b-cdn.net/path/to/image.jpg/w_299/crop_smart`
