import React, { useEffect, useState } from "react";
import ArrowForwardIcon from "@mui/icons-material/ArrowForward";
import FileUploadRoundedIcon from "@mui/icons-material/FileUploadRounded";
import IconButton from "@mui/material/IconButton";
import CircularProgress from "@mui/material/CircularProgress";
import InputBase from "@mui/material/InputBase";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import { message } from "antd";
import axios from "axios";
import PhotoAlbum from "react-photo-album";
import Lightbox from "yet-another-react-lightbox";
import Fullscreen from "yet-another-react-lightbox/plugins/fullscreen";
import Slideshow from "yet-another-react-lightbox/plugins/slideshow";
import Thumbnails from "yet-another-react-lightbox/plugins/thumbnails";
import Zoom from "yet-another-react-lightbox/plugins/zoom";
import styled from "styled-components";
import { BASE_URL, TOKEN_KEY } from "../constants";

const Overlay = styled.div`
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
`;

const MainContainer = styled.div`
  background-color: #27272a;
  min-height: 100vh;
  padding: 112px 24px 48px;
`;

const HeaderContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

function Landing() {
  const [index, setIndex] = useState(-1);
  const [inputValue, setInputValue] = useState("");
  const [isGeneratingImage, setIsGeneratingImage] = useState(false);
  const [generatedImageUrl, setGeneratedImageUrl] = useState("");
  const [slicedPhotos, setSlicedPhotos] = useState([]);

  useEffect(() => {
    if (generatedImageUrl) {
      setSlicedPhotos([{ src: generatedImageUrl, width: 200, height: 200 }]);
    }
  }, [generatedImageUrl]);

  const createImage = async () => {
    if (!inputValue.trim()) {
      message.error("Enter an image description first.");
      return;
    }

    try {
      setIsGeneratingImage(true);
      const response = await axios({
        method: "POST",
        url: `${BASE_URL}/generate-image`,
        headers: {
          Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`,
          "Content-Type": "application/json",
        },
        data: {
          prompt: inputValue.trim(),
        },
      });
      const imageDataUrl = response.data?.imageDataUrl;
      if (!imageDataUrl) {
        throw new Error("The image API did not return image data.");
      }
      setGeneratedImageUrl(imageDataUrl);
    } catch (error) {
      console.error("image generation failed:", error);
      message.error(error.response?.data || "Image generation failed.");
    } finally {
      setIsGeneratingImage(false);
    }
  };

  const handleUploadImage = async () => {
    if (!generatedImageUrl) {
      return;
    }

    try {
      const response = await fetch(generatedImageUrl);
      const blob = await response.blob();
      const formData = new FormData();
      formData.append("message", `AI generated image: ${inputValue}`);
      formData.append("media_file", new File([blob], "ai-generated.png", { type: "image/png" }));

      await axios({
        method: "POST",
        url: `${BASE_URL}/upload`,
        headers: {
          Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`,
        },
        data: formData,
      });
      message.success("AI-generated image uploaded.");
      setIndex(-1);
    } catch (error) {
      console.error("AI image upload failed:", error);
      message.error("Upload failed.");
    }
  };

  return (
    <MainContainer>
      {isGeneratingImage && (
        <Overlay>
          <CircularProgress color="info" size={80} />
        </Overlay>
      )}
      <HeaderContainer>
        <Typography variant="h1" className="landing-title">
          Social AI
        </Typography>
        <Typography variant="h5" className="landing-subtitle">
          Unleash creativity and share memories where AI meets your imagination.
        </Typography>
        <Paper component="form" className="prompt-box" onSubmit={(event) => event.preventDefault()}>
          <InputBase
            multiline
            sx={{ ml: 1, flex: 1 }}
            placeholder="Enter your image description here"
            value={inputValue}
            onChange={(event) => setInputValue(event.target.value)}
          />
          <IconButton type="button" onClick={createImage} aria-label="generate image">
            <ArrowForwardIcon />
          </IconButton>
        </Paper>
      </HeaderContainer>

      {slicedPhotos.length > 0 && (
        <PhotoAlbum
          photos={slicedPhotos}
          layout="rows"
          targetRowHeight={200}
          onClick={({ index: photoIndex }) => setIndex(photoIndex)}
        />
      )}
      <Lightbox
        slides={slicedPhotos}
        open={index >= 0}
        index={index}
        close={() => setIndex(-1)}
        plugins={[Fullscreen, Slideshow, Thumbnails, Zoom]}
        toolbar={{
          buttons: [
            <IconButton
              key="upload"
              type="button"
              aria-label="upload generated image"
              onClick={handleUploadImage}
            >
              <FileUploadRoundedIcon sx={{ color: "white" }} />
            </IconButton>,
          ],
        }}
      />
    </MainContainer>
  );
}

export default Landing;
