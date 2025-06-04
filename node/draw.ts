import { createCanvas, type CanvasRenderingContext2D } from "canvas";

function drawTriangle(
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  size: number,
  maxDepth: number
) {
  const queue: Array<{ x: number; y: number; size: number; depth: number }> = [
    { x, y, size, depth: 0 },
  ];

  while (queue.length > 0) {
    const { x, y, size, depth } = queue.shift()!;

    if (depth === maxDepth || size < 1) {
      ctx.beginPath();
      ctx.moveTo(x, y);
      ctx.lineTo(x + size / 2, y + (Math.sqrt(3) * size) / 2);
      ctx.lineTo(x - size / 2, y + (Math.sqrt(3) * size) / 2);
      ctx.closePath();
      ctx.fillStyle = "black";
      ctx.fill();
    } else {
      const half = size / 2;
      const height = (Math.sqrt(3) * half) / 2;

      queue.push({ x, y, size: half, depth: depth + 1 });
      queue.push({
        x: x - half / 2,
        y: y + height,
        size: half,
        depth: depth + 1,
      });
      queue.push({
        x: x + half / 2,
        y: y + height,
        size: half,
        depth: depth + 1,
      });
    }
  }
}

export const createImage = (depth = 8) => {
  console.log("Tiefe:", depth);
  console.log("Start Zeit", new Date().toISOString());
  console.time("createImage");

  const size = 20000;
  const canvas = createCanvas(size, size);
  const ctx = canvas.getContext("2d");
  ctx.fillStyle = "white";
  ctx.fillRect(0, 0, size, size);
  const centerX = size / 2;
  const topY = 0;
  drawTriangle(ctx, centerX, topY, size * 0.9, depth);

  console.log("End Zeit:", new Date().toISOString());
  console.timeEnd("createImage");
  return canvas;
};
