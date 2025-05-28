import express from 'express';
import { createImage } from './draw';
import type { Request, Response } from 'express';

const app = express();

app.get('/api/status', (_req: Request, res: Response) => {
  res.send('OK');
});

app.get('/api/image/:depth', async (req: Request, res: Response) => {
  try {
    const depth = req.params.depth;
    const img = createImage(parseInt(depth));
    const buffer = await img.toBuffer('image/png');
      res.setHeader('Content-Type', 'image/png');
      res.setHeader('Cache-Control', 'no-store, no-cache, must-revalidate');
      res.setHeader('Pragma', 'no-cache');
      res.setHeader('Expires', '0');
      res.setHeader('Content-Disposition', 'inline; filename="image.png"');
      res.setHeader('Content-Type', 'image/png');
    res.send(buffer);
  } catch (err) {
    console.error('Image generation error:', err);
    res.status(500).send('Internal Server Error');
  }
});

app.use((_req: Request, res: Response) => {
  res.status(404).send('Not Found');
});

app.listen(3001);
