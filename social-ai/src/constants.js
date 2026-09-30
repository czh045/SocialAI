export const TOKEN_KEY = "token";

// The lesson used a deployed GAE URL here. Local development uses the Go
// backend on port 8080, while production can set REACT_APP_API_BASE_URL.
export const BASE_URL =
  process.env.REACT_APP_API_BASE_URL || "http://localhost:8080";

export const SEARCH_KEY = {
  all: 0,
  keyword: 1,
  user: 2,
};
