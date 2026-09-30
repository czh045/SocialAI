import React, { useEffect, useState } from "react";
import { Col, message, Row, Tabs } from "antd";
import axios from "axios";
import { BASE_URL, SEARCH_KEY, TOKEN_KEY } from "../constants";
import CreatePostButton from "./CreatePostButton";
import PhotoGallery from "./PhotoGallery";
import SearchBar from "./SearchBar";

const { TabPane } = Tabs;

function Collection() {
  const [posts, setPosts] = useState([]);
  const [activeTab, setActiveTab] = useState("image");
  const [searchOption, setSearchOption] = useState({
    type: SEARCH_KEY.all,
    keyword: "",
  });

  useEffect(() => {
    const { type, keyword } = searchOption;
    let url = `${BASE_URL}/search`;
    if (type === SEARCH_KEY.user) {
      url = `${BASE_URL}/search?user=${encodeURIComponent(keyword)}`;
    } else if (type === SEARCH_KEY.keyword) {
      url = `${BASE_URL}/search?keywords=${encodeURIComponent(keyword)}`;
    }

    axios({
      method: "GET",
      url,
      headers: {
        Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`,
      },
    })
      .then((response) => {
        if (response.status === 200) {
          setPosts(response.data);
        }
      })
      .catch((error) => {
        console.error("fetch posts failed:", error.message);
        message.error("Fetch posts failed.");
      });
  }, [searchOption]);

  const renderPosts = (type) => {
    const filtered = posts.filter((post) => post.type === type);
    if (filtered.length === 0) {
      return <div className="empty-result">No {type === "image" ? "images" : "videos"} yet.</div>;
    }

    if (type === "image") {
      const imageArr = filtered.map((image) => ({
        postId: image.id,
        src: image.url,
        user: image.user,
        caption: image.message,
      }));
      return <PhotoGallery images={imageArr} />;
    }

    return (
      <Row gutter={[16, 16]}>
        {filtered.map((post) => (
          <Col xs={24} md={12} key={post.id}>
            <video src={post.url} controls className="video-block" />
            <div className="video-caption">{`${post.user}: ${post.message}`}</div>
          </Col>
        ))}
      </Row>
    );
  };

  const showPost = (type) => {
    setActiveTab(type);
    setTimeout(() => {
      setSearchOption({ type: SEARCH_KEY.all, keyword: "" });
    }, 300);
  };

  return (
    <div className="home">
      <SearchBar handleSearch={setSearchOption} />
      <div className="display">
        <Tabs
          onChange={setActiveTab}
          defaultActiveKey="image"
          activeKey={activeTab}
          tabBarExtraContent={<CreatePostButton onShowPost={showPost} />}
        >
          <TabPane tab="Images" key="image">
            {renderPosts("image")}
          </TabPane>
          <TabPane tab="Videos" key="video">
            {renderPosts("video")}
          </TabPane>
        </Tabs>
      </div>
    </div>
  );
}

export default Collection;
