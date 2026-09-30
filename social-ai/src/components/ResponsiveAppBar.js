import React, { useState } from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import IconButton from "@mui/material/IconButton";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import LogoutIcon from "@mui/icons-material/Logout";
import MenuIcon from "@mui/icons-material/Menu";
import PublicIcon from "@mui/icons-material/Public";
import { useNavigate } from "react-router-dom";

const pages = ["Create", "Collection"];

function ResponsiveAppBar({ isLoggedIn, handleLogout }) {
  const [anchorElNav, setAnchorElNav] = useState(null);
  const navigate = useNavigate();

  const handleOpenNavMenu = (event) => setAnchorElNav(event.currentTarget);
  const handleCloseNavMenu = () => setAnchorElNav(null);
  const handleNavigate = (page) => {
    navigate(`/${page.toLowerCase()}`);
    handleCloseNavMenu();
  };

  return (
    <AppBar
      position="fixed"
      sx={{ top: 0, zIndex: 10, backgroundColor: "#000000" }}
    >
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          <Button
            color="inherit"
            onClick={() => navigate(isLoggedIn ? "/create" : "/login")}
            sx={{ minWidth: 0, mr: 1, p: 0.5 }}
            aria-label="Social AI home"
          >
            <PublicIcon />
          </Button>
          <Typography
            variant="h6"
            component="button"
            onClick={() => navigate(isLoggedIn ? "/create" : "/login")}
            sx={{
              mr: 2,
              border: 0,
              background: "transparent",
              color: "inherit",
              cursor: "pointer",
              fontFamily: "monospace",
              fontWeight: 700,
              letterSpacing: "0.18rem",
            }}
          >
            Social AI
          </Typography>

          {isLoggedIn && (
            <>
              <Box sx={{ flexGrow: 1, display: { xs: "flex", md: "none" } }}>
                <IconButton
                  size="large"
                  aria-label="open navigation menu"
                  onClick={handleOpenNavMenu}
                  color="inherit"
                >
                  <MenuIcon />
                </IconButton>
                <Menu
                  anchorEl={anchorElNav}
                  anchorOrigin={{ vertical: "bottom", horizontal: "left" }}
                  keepMounted
                  transformOrigin={{ vertical: "top", horizontal: "left" }}
                  open={Boolean(anchorElNav)}
                  onClose={handleCloseNavMenu}
                  sx={{ display: { xs: "block", md: "none" } }}
                >
                  {pages.map((page) => (
                    <MenuItem key={page} onClick={() => handleNavigate(page)}>
                      <Typography textAlign="center">{page}</Typography>
                    </MenuItem>
                  ))}
                </Menu>
              </Box>
              <Box sx={{ flexGrow: 1, display: { xs: "none", md: "flex" } }}>
                {pages.map((page) => (
                  <Button
                    key={page}
                    onClick={() => handleNavigate(page)}
                    sx={{ my: 1, color: "white" }}
                  >
                    {page}
                  </Button>
                ))}
              </Box>
              <Button
                variant="outlined"
                startIcon={<LogoutIcon />}
                onClick={handleLogout}
                sx={{ backgroundColor: "white", color: "black", borderColor: "white" }}
              >
                Log out
              </Button>
            </>
          )}
        </Toolbar>
      </Container>
    </AppBar>
  );
}

export default ResponsiveAppBar;
