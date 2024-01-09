import base64
from io import BytesIO
from typing import Literal

import mandelbrot
from fastapi import FastAPI
from fastapi.responses import PlainTextResponse, RedirectResponse, Response
from pydantic import BaseModel


class Request(BaseModel):
    width: int = 640
    height: int = 480
    iterations: int = 100
    re_min: float = -2.0
    re_max: float = 1.0
    im_min: float = -1.0
    im_max: float = 1.0
    kind: Literal["png", "base64"] = "png"


app = FastAPI(title="Mandelbrot")


@app.get("/")
def index():
    return RedirectResponse('/docs')


@app.post("/generate/")
def generate_image(req: Request):
    img = mandelbrot.generate(
        req.width,
        req.height,
        req.iterations,
        req.re_min,
        req.re_max,
        req.im_min,
        req.im_max,
        req.kind,
    )
    buffer = BytesIO()
    img.save(buffer, format="png")
    payload = buffer.getbuffer().tobytes()
    if req.kind == 'png':
        return Response(content=payload, media_type="image/png")
    else:
        return PlainTextResponse(content=base64.b64encode(payload))
