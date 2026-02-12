"use client";

import { FormEvent, useMemo, useState } from "react";
import styles from "./page.module.css";

type ApiResponse<T = unknown> = {
  success: boolean;
  message: string;
  data?: T;
};

type LoginResponse = {
  accessToken: string;
  user: {
    id: string;
    name: string;
    email: string;
    role: "teacher" | "student";
  };
};

export default function Home() {
  const [token, setToken] = useState("");
  const [output, setOutput] = useState("Ready");

  const [register, setRegister] = useState({
    name: "",
    email: "",
    password: "",
  });

  const [login, setLogin] = useState({
    email: "",
    password: "",
  });

  const [className, setClassName] = useState("");
  const [classId, setClassId] = useState("");
  const [studentId, setStudentId] = useState("");

  const authHeader = useMemo(
    () => (token ? { Authorization: `Bearer ${token}` } : {}),
    [token],
  );

  async function apiCall<T>(
    path: string,
    method: string,
    body?: unknown,
    useAuth = false,
  ) {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };
    if (useAuth && authHeader.Authorization) {
      headers.Authorization = authHeader.Authorization;
    }

    const response = await fetch(`/api/proxy/${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    const result = (await response.json()) as ApiResponse<T>;
    setOutput(JSON.stringify(result, null, 2));
    return result;
  }

  async function onRegister(e: FormEvent) {
    e.preventDefault();
    await apiCall("users/register", "POST", register);
  }

  async function onLogin(e: FormEvent) {
    e.preventDefault();
    const result = await apiCall<LoginResponse>("users/login", "POST", login);
    const accessToken = (result.data as LoginResponse | undefined)?.accessToken;
    if (accessToken) {
      setToken(accessToken);
    }
  }

  async function onMe() {
    await apiCall("users/me", "GET", undefined, true);
  }

  async function onCreateClass(e: FormEvent) {
    e.preventDefault();
    const result = await apiCall<{ id?: string }>(
      "class",
      "POST",
      { class_name: className },
      true,
    );
    const createdId = (result.data as { id?: string } | undefined)?.id;
    if (createdId) {
      setClassId(createdId);
    }
  }

  async function onGetClass() {
    const normalizedClassId = classId.trim().replace(/^\/+|\/+$/g, "");
    if (!normalizedClassId) return;
    await apiCall(`class/${normalizedClassId}`, "GET", undefined, true);
  }

  async function onAddStudent(e: FormEvent) {
    e.preventDefault();
    const normalizedClassId = classId.trim().replace(/^\/+|\/+$/g, "");
    const normalizedStudentId = studentId.trim().replace(/^\/+|\/+$/g, "");
    if (!normalizedClassId || !normalizedStudentId) return;

    await apiCall(
      "class",
      "PATCH",
      { class_id: normalizedClassId, student_id: normalizedStudentId },
      true,
    );
  }

  async function onStartAttendance(e: FormEvent) {
    e.preventDefault();
    const normalizedClassId = classId.trim().replace(/^\/+|\/+$/g, "");
    if (!normalizedClassId) return;
    await apiCall("attendance/start", "POST", { class_id: normalizedClassId }, true);
  }

  async function onMyAttendance() {
    const normalizedClassId = classId.trim().replace(/^\/+|\/+$/g, "");
    if (!normalizedClassId) return;

    const queryStudentId = studentId.trim();
    const path = queryStudentId
      ? `class/${normalizedClassId}/my-attendance?student_id=${encodeURIComponent(queryStudentId)}`
      : `class/${normalizedClassId}/my-attendance`;

    await apiCall(path, "GET", undefined, true);
  }

  return (
    <main className={styles.page}>
      <section className={styles.hero}>
        <h1>Attendance Dashboard</h1>
        <p>Basic Next.js client for your Go attendance API.</p>
      </section>

      <section className={styles.grid}>
        <article className={styles.card}>
          <h2>Register (student)</h2>
          <form onSubmit={onRegister} className={styles.form}>
            <input
              placeholder="Name"
              value={register.name}
              onChange={(e) => setRegister((v) => ({ ...v, name: e.target.value }))}
            />
            <input
              placeholder="Email"
              value={register.email}
              onChange={(e) =>
                setRegister((v) => ({ ...v, email: e.target.value }))
              }
            />
            <input
              placeholder="Password"
              type="password"
              value={register.password}
              onChange={(e) =>
                setRegister((v) => ({ ...v, password: e.target.value }))
              }
            />
            <button type="submit">Register</button>
          </form>
        </article>

        <article className={styles.card}>
          <h2>Login</h2>
          <form onSubmit={onLogin} className={styles.form}>
            <input
              placeholder="Email"
              value={login.email}
              onChange={(e) => setLogin((v) => ({ ...v, email: e.target.value }))}
            />
            <input
              placeholder="Password"
              type="password"
              value={login.password}
              onChange={(e) =>
                setLogin((v) => ({ ...v, password: e.target.value }))
              }
            />
            <button type="submit">Login</button>
          </form>
          <p className={styles.small}>Token: {token ? "Loaded" : "Missing"}</p>
          <button onClick={onMe}>Get /users/me</button>
        </article>

        <article className={styles.card}>
          <h2>Class + Attendance</h2>
          <form onSubmit={onCreateClass} className={styles.form}>
            <input
              placeholder="Class Name"
              value={className}
              onChange={(e) => setClassName(e.target.value)}
            />
            <button type="submit">Create Class (teacher)</button>
          </form>

          <div className={styles.form}>
            <input
              placeholder="Class ID"
              value={classId}
              onChange={(e) => setClassId(e.target.value)}
            />
            <button onClick={onGetClass}>Get Class</button>
            <button onClick={onMyAttendance}>My Attendance (student)</button>
          </div>

          <form onSubmit={onAddStudent} className={styles.form}>
            <input
              placeholder="Student ID"
              value={studentId}
              onChange={(e) => setStudentId(e.target.value)}
            />
            <button type="submit">Add Student (teacher)</button>
          </form>

          <form onSubmit={onStartAttendance} className={styles.form}>
            <button type="submit">Start Attendance (teacher)</button>
          </form>
        </article>
      </section>

      <section className={styles.output}>
        <h2>API Response</h2>
        <pre>{output}</pre>
      </section>
    </main>
  );
}
