import { createImage } from "./draw.ts";

Deno.serve({ port: 3002 }, async (req) => {
  const url = new URL(req.url);
  const pathname = url.pathname;

  if (pathname === "/api/status") {
    return new Response("OK");
  }

  const match = pathname.match(/^\/api\/image\/(\d+)$/);
  if (match) {
    const depth = parseInt(match[1]);

    try {
      const img = createImage(depth);
      const buffer = await img.toBuffer("image/png");

      return new Response(buffer, {
      headers: {
        "Content-Type": "image/png",
        "Cache-Control": "no-store, no-cache, must-revalidate",
        "Pragma": "no-cache",
        "Expires": "0",
        "Content-Disposition": `inline; filename="image.png"`,
      },
      });
    } catch (err) {
      console.error("Error generating image:", err);
      return new Response("Internal Server Error", { status: 500 });
    }
  }

  return new Response("Not Found", { status: 404 });
});
