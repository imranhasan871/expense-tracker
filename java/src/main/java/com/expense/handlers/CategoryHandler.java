package com.expense.handlers;

import com.expense.config.DatabaseConfig;
import com.expense.models.Category;
import com.expense.repository.CategoryRepository;
import com.sun.net.httpserver.HttpExchange;
import java.io.IOException;
import java.sql.Connection;
import java.sql.SQLException;
import java.util.List;
import java.io.InputStreamReader;
import java.io.BufferedReader;

public class CategoryHandler extends BaseHandler {

  @Override
  public void handle(HttpExchange exchange) throws IOException {
    String method = exchange.getRequestMethod();

    try (Connection conn = DatabaseConfig.getConnection()) {
      CategoryRepository repo = new CategoryRepository(conn);

      if ("GET".equalsIgnoreCase(method)) {
        List<Category> categories = repo.findAll();
        sendResponse(exchange, 200, new SuccessResponse(categories));
      } else if ("POST".equalsIgnoreCase(method)) {
        // Read Request Body
        Category req = gson.fromJson(new InputStreamReader(exchange.getRequestBody()), Category.class);
        repo.create(req);
        sendResponse(exchange, 201, new SuccessResponse(null, "Category created"));
      } else {
        sendError(exchange, 405, "Method Not Allowed");
      }
    } catch (SQLException e) {
      e.printStackTrace();
      sendError(exchange, 500, e.getMessage());
    }
  }
}
