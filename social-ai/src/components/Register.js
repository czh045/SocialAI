import React from "react";
import { useNavigate } from "react-router-dom";
import { Button, Form, Input, message } from "antd";
import axios from "axios";
import { BASE_URL } from "../constants";

const formItemLayout = {
  labelCol: { xs: { span: 24 }, sm: { span: 8 } },
  wrapperCol: { xs: { span: 24 }, sm: { span: 16 } },
};

const tailFormItemLayout = {
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 16, offset: 8 },
  },
};

function Register() {
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const onFinish = ({ username, password }) => {
    axios({
      method: "POST",
      url: `${BASE_URL}/signup`,
      data: { username, password },
      headers: { "Content-Type": "application/json" },
    })
      .then((response) => {
        if (response.status === 200) {
          message.success("Registration succeeded!");
          navigate("/login");
        }
      })
      .catch((error) => {
        console.error("registration failed:", error.message);
        message.error("Registration failed. Try a different username.");
      });
  };

  return (
    <Form
      {...formItemLayout}
      form={form}
      name="register"
      onFinish={onFinish}
      className="register"
    >
      <Form.Item
        name="username"
        label="Username"
        rules={[{ required: true, message: "Please input your username." }]}
      >
        <Input />
      </Form.Item>
      <Form.Item
        name="password"
        label="Password"
        rules={[{ required: true, message: "Please input your password." }]}
        hasFeedback
      >
        <Input.Password />
      </Form.Item>
      <Form.Item
        name="confirm"
        label="Confirm password"
        dependencies={["password"]}
        hasFeedback
        rules={[
          { required: true, message: "Please confirm your password." },
          ({ getFieldValue }) => ({
            validator(_, value) {
              if (!value || getFieldValue("password") === value) {
                return Promise.resolve();
              }
              return Promise.reject(
                new Error("The two passwords that you entered do not match.")
              );
            },
          }),
        ]}
      >
        <Input.Password />
      </Form.Item>
      <Form.Item {...tailFormItemLayout}>
        <Button type="primary" htmlType="submit" className="register-btn">
          Register
        </Button>
      </Form.Item>
    </Form>
  );
}

export default Register;
