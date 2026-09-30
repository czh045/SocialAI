import React from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import Collection from "./Collection";
import Landing from "./Landing";
import Login from "./Login";
import Register from "./Register";

function Main({ isLoggedIn, handleLoggedIn }) {
  const showLogin = () =>
    isLoggedIn ? <Navigate to="/create" replace /> : <Login handleLoggedIn={handleLoggedIn} />;

  const showRegister = () =>
    isLoggedIn ? <Navigate to="/create" replace /> : <Register />;

  const showLanding = () =>
    isLoggedIn ? <Landing /> : <Navigate to="/login" replace />;

  const showCollection = () =>
    isLoggedIn ? <Collection /> : <Navigate to="/login" replace />;

  return (
    <main className="main">
      <Routes>
        <Route path="/" element={showLogin()} />
        <Route path="/login" element={showLogin()} />
        <Route path="/register" element={showRegister()} />
        <Route path="/create" element={showLanding()} />
        <Route path="/collection" element={showCollection()} />
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </main>
  );
}

export default Main;
