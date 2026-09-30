import React from "react";
import { Button, Form, Input, message } from "antd";
import { LockOutlined, UserOutlined } from "@ant-design/icons";
import axios from "axios";
import { Link } from "react-router-dom";
import { BASE_URL } from "../constants";

function Login({ handleLoggedIn }) {
  const onFinish = ({ username, password }) => {
    axios({
      method: "POST",
      url: `${BASE_URL}/signin`,
      data: { username, password },
      headers: { "Content-Type": "application/json" },
    })
      .then((response) => {
        if (response.status === 200) {
          handleLoggedIn(response.data);
          message.success("Login succeeded!");
        }
      })
      .catch((error) => {
        console.error("login failed:", error.message);
        message.error("Login failed. Check your username and password.");
      });
  };

  return (
    <Form name="normal_login" className="login-form" onFinish={onFinish}>
      <Form.Item
        name="username"
        rules={[{ required: true, message: "Please input your username." }]}
      >
        <Input prefix={<UserOutlined />} placeholder="Username" />
      </Form.Item>
      <Form.Item
        name="password"
        rules={[{ required: true, message: "Please input your password." }]}
      >
        <Input prefix={<LockOutlined />} type="password" placeholder="Password" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" className="login-form-button">
          Log in
        </Button>
        <span className="form-link">
          Or <Link to="/register">register now</Link>
        </span>
      </Form.Item>
    </Form>
  );
}

export default Login;
