# Next Steps for Session Management

This document outlines the steps to implement session management for the application.

## 1. Implement Server-Side Session Store
   - Define a `Session` struct (e.g., `ID`, `UserID`, `Expiry`).
   - Create an in-memory map or a database table to store active sessions.
   - Implement functions to:
     - `CreateSession(userID int)`: Generates a new session ID, stores the session, and returns the session ID.
     - `GetSession(sessionID string)`: Retrieves a session by its ID.
     - `DeleteSession(sessionID string)`: Removes a session from the store.

## 2. Modify LoginHandler to Create Sessions
   - After successful user authentication in `LoginHandler`:
     - Call `CreateSession` to generate a new session for the logged-in user.
     - Set an `http.Cookie` in the response containing the session ID. This cookie should be `HttpOnly` and `Secure`.

## 3. Implement Session Validation Middleware
   - Create an HTTP middleware function (e.g., `AuthMiddleware`) that:
     - Checks for the session cookie in incoming requests.
     - Retrieves the session ID from the cookie.
     - Uses `GetSession` to validate the session ID and check its expiry.
     - If valid, attach the `UserID` to the request context.
     - If invalid or missing, redirect the user to the login page or return an unauthorized error.

## 4. Protect Dashboard and Other Routes
   - Apply the `AuthMiddleware` to the `/dashboard.html` route and any other routes that require authentication.

## 5. Implement Logout Functionality
   - Create a new API endpoint (e.g., `/api/logout`).
   - In the `LogoutHandler`:
     - Delete the session from the server-side store using `DeleteSession`.
     - Clear the session cookie from the client's browser.
     - Redirect the user to the login page.
   - Add JavaScript to `dashboard.html` to call this logout endpoint when the "Sign Out" button is clicked.

## 6. Update Client-Side (index.html and dashboard.html)
   - In `index.html`: Store the username (returned from `LoginHandler`) in `localStorage` upon successful login.
   - In `dashboard.html`: Retrieve the username from `localStorage` and display it.
   - Add JavaScript for "Sign Out" and "Delete Account" buttons (initially just for logout).
