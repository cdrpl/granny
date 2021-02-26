import express from "express";
import morgan from "morgan";
import type { ErrorRequestHandler, Express, Router } from "express";

function createApp(router: Router): Express {
  const app = express();
  app.disable("x-powered-by");
  app.use(express.urlencoded({ extended: true }));

  // Dev logging
  if (process.env.NODE_ENV === "development") {
    app.use(morgan("dev"));
  }

  app.use(router);

  // Catch 404
  app.use((req, res) => {
    res.status(404).json({ error: { code: 404, message: "Not Found" } });
  });

  // Error handler
  const errHandler: ErrorRequestHandler = (err, req, res, next) => {
    let { message } = err;

    console.error(err);

    if (process.env.NODE_ENV === "production") {
      message = "An error has occured"; // Hide err message in production
    }

    res.status(500).json({ error: { code: 500, message } });
  };

  app.use(errHandler);

  return app;
}

export { createApp };
