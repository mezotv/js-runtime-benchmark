import { createImage } from "./draw.ts";

Bun.serve({

  routes: {
    // Static routes
    "/api/status": new Response("OK"),

    "/api/image/:depth": async (req: { params: { depth: string; }; }) => {
      const depth = req.params.depth;
        const img = createImage(parseInt(depth));
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
    }
  },


  fetch(req) {
    return new Response("Not Found", { status: 404 });
  },
});