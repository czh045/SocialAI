import React, { forwardRef } from "react";
import { Form, Input, Upload } from "antd";
import { InboxOutlined } from "@ant-design/icons";

export const PostForm = forwardRef((_, formRef) => {
  const formItemLayout = {
    labelCol: { span: 6 },
    wrapperCol: { span: 14 },
  };

  const normalizeFile = (event) =>
    Array.isArray(event) ? event : event?.fileList;

  return (
    <Form name="post-form" {...formItemLayout} ref={formRef}>
      <Form.Item
        name="description"
        label="Message"
        rules={[{ required: true, message: "Please input your message." }]}
      >
        <Input />
      </Form.Item>
      <Form.Item label="Media">
        <Form.Item
          name="uploadPost"
          valuePropName="fileList"
          getValueFromEvent={normalizeFile}
          noStyle
          rules={[{ required: true, message: "Please select an image or video." }]}
        >
          <Upload.Dragger name="files" beforeUpload={() => false} maxCount={1}>
            <p className="ant-upload-drag-icon">
              <InboxOutlined />
            </p>
            <p className="ant-upload-text">
              Click or drag one image/video file to this area
            </p>
          </Upload.Dragger>
        </Form.Item>
      </Form.Item>
    </Form>
  );
});
