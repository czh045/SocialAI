import React, { useEffect, useMemo, useState } from "react";
import PropTypes from "prop-types";
import { Button, message } from "antd";
import { DeleteOutlined } from "@ant-design/icons";
import DeleteIcon from "@mui/icons-material/Delete";
import IconButton from "@mui/material/IconButton";
import axios from "axios";
import PhotoAlbum from "react-photo-album";
import Lightbox from "yet-another-react-lightbox";
import Fullscreen from "yet-another-react-lightbox/plugins/fullscreen";
import Slideshow from "yet-another-react-lightbox/plugins/slideshow";
import Thumbnails from "yet-another-react-lightbox/plugins/thumbnails";
import Zoom from "yet-another-react-lightbox/plugins/zoom";
import { BASE_URL, TOKEN_KEY } from "../constants";

const captionStyle = {
  backgroundColor: "rgba(0, 0, 0, 0.6)",
  position: "absolute",
  bottom: 0,
  width: "100%",
  color: "white",
  padding: "4px",
  fontSize: "90%",
};

function PhotoGallery({ images }) {
  const [galleryImages, setGalleryImages] = useState(images);
  const [index, setIndex] = useState(-1);

  useEffect(() => setGalleryImages(images), [images]);

  const deleteImage = (postId) => {
    if (!postId || !window.confirm("Are you sure you want to delete this image?")) {
      return;
    }

    axios({
      method: "DELETE",
      url: `${BASE_URL}/post/${postId}`,
      headers: {
        Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`,
      },
    })
      .then((response) => {
        if (response.status === 200) {
          setGalleryImages((current) =>
            current.filter((image) => image.postId !== postId)
          );
          setIndex(-1);
        }
      })
      .catch((error) => {
        console.error("delete post failed:", error.message);
        message.error("Delete post failed.");
      });
  };

  const imageArr = useMemo(
    () =>
      galleryImages.map((image) => ({
        ...image,
        width: 300,
        height: 200,
        title: `${image.user}: ${image.caption}`,
      })),
    [galleryImages]
  );

  return (
    <div className="gallery-wrapper">
      <PhotoAlbum
        photos={imageArr}
        layout="rows"
        targetRowHeight={200}
        onClick={({ index: photoIndex }) => setIndex(photoIndex)}
        renderPhoto={({ imageProps, photo }) => (
          <div className="gallery-photo">
            <img {...imageProps} alt={photo.caption || "Social AI post"} />
            <div style={captionStyle}>
              <div>{`${photo.user}: ${photo.caption}`}</div>
              <Button
                size="small"
                type="primary"
                icon={<DeleteOutlined />}
                onClick={() => deleteImage(photo.postId)}
              >
                Delete
              </Button>
            </div>
          </div>
        )}
      />
      <Lightbox
        slides={imageArr}
        open={index >= 0}
        index={index}
        close={() => setIndex(-1)}
        plugins={[Fullscreen, Slideshow, Thumbnails, Zoom]}
        toolbar={{
          buttons: [
            <IconButton
              key="delete"
              aria-label="delete selected image"
              onClick={() => deleteImage(imageArr[index]?.postId)}
            >
              <DeleteIcon sx={{ color: "#cccccc" }} />
            </IconButton>,
          ],
        }}
      />
    </div>
  );
}

PhotoGallery.propTypes = {
  images: PropTypes.arrayOf(
    PropTypes.shape({
      postId: PropTypes.string.isRequired,
      user: PropTypes.string.isRequired,
      caption: PropTypes.string.isRequired,
      src: PropTypes.string.isRequired,
    })
  ).isRequired,
};

export default PhotoGallery;
