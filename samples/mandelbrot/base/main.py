import asyncio
import base64
import logging
import mandelbrot
import uvicorn
from fastapi import FastAPI, status
from fastapi.responses import RedirectResponse
from io import BytesIO
from pydantic import BaseModel


class EndpointFilter(logging.Filter):
    def filter(self, record: logging.LogRecord) -> bool:
        return record.args and len(record.args) >= 3 and record.args[2] != "/health"


# Filter out noisy health check logs
logging.getLogger("uvicorn.access").addFilter(EndpointFilter())

class HealthResponse(BaseModel):
    status: str = "OK"


class GenerateImageRequest(BaseModel):
    width: int = 640
    height: int = 480
    iterations: int = 100
    re_min: float = -2.0
    re_max: float = 1.0
    im_min: float = -1.0
    im_max: float = 1.0
    delay: int = 0


class GenerateImageResponse(BaseModel):
    image: str


app = FastAPI(title="Mandelbrot")


@app.get("/")
def index():
    return RedirectResponse("/docs")


@app.post(
        "/generate",
        tags=["mandelbrot"],
        summary="Generate an image of a Mandelbrot set",
        response_description="Returns HTTP status 200 OK with the generated image",
        status_code=status.HTTP_200_OK,
        response_model=GenerateImageResponse)
async def generate_image(req: GenerateImageRequest):
    if req.delay != 0:
        await asyncio.sleep(req.delay)
    img = mandelbrot.generate(
        req.width,
        req.height,
        req.iterations,
        req.re_min,
        req.re_max,
        req.im_min,
        req.im_max,
    )
    buffered = BytesIO()
    img.save(buffered, format="png")
    return GenerateImageResponse(image=base64.b64encode(buffered.getvalue()).decode("utf-8"))


@app.get(
        "/health",
        tags=["healthcheck"],
        summary="Perform a Health Check",
        response_description="Returns HTTP status 200 OK",
        status_code=status.HTTP_200_OK,
        response_model=HealthResponse)
def get_health():
    return HealthResponse(status="OK")
