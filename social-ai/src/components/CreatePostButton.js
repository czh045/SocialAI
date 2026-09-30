import React, { Component } from "react";
import { Button, message, Modal } from "antd";
import axios from "axios";
import { BASE_URL, TOKEN_KEY } from "../constants";
import { PostForm } from "./PostForm";

class CreatePostButton extends Component {
  state = {
    visible: false,
    confirmLoading: false,
  };

  showModal = () => this.setState({ visible: true });

  handleOk = () => {
    this.setState({ confirmLoading: true });

    this.postForm
      .validateFields()
      .then(({ description, uploadPost }) => {
        const { type, originFileObj } = uploadPost[0];
        const postType = type?.match(/^(image|video)/)?.[0];

        if (!postType) {
          message.error("Please choose a valid image or video.");
          this.setState({ confirmLoading: false });
          return;
        }

        const formData = new FormData();
        formData.append("message", description);
        formData.append("media_file", originFileObj);

        return axios({
          method: "POST",
          url: `${BASE_URL}/upload`,
          headers: {
            Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`,
          },
          data: formData,
        })
          .then((response) => {
            if (response.status === 200) {
              message.success("The image/video is uploaded!");
              this.postForm.resetFields();
              this.handleCancel();
              this.props.onShowPost(postType);
            }
          })
          .catch((error) => {
            console.error("upload failed:", error.message);
            message.error("Failed to upload image/video.");
          })
          .finally(() => this.setState({ confirmLoading: false }));
      })
      .catch(() => this.setState({ confirmLoading: false }));
  };

  handleCancel = () => this.setState({ visible: false });

  render() {
    const { visible, confirmLoading } = this.state;

    return (
      <div>
        <Button type="primary" onClick={this.showModal}>
          Create New Post
        </Button>
        <Modal
          title="Create New Post"
          open={visible}
          onOk={this.handleOk}
          okText="Create"
          confirmLoading={confirmLoading}
          onCancel={this.handleCancel}
        >
          <PostForm ref={(instance) => (this.postForm = instance)} />
        </Modal>
      </div>
    );
  }
}

export default CreatePostButton;
