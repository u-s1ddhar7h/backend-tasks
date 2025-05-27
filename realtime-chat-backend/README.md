# 💬 Real-Time Chat Backend

This is a backend server for a real-time chat application built using **Node.js** and **Socket.IO**. It supports multiple chat rooms and enables real-time messaging between users via WebSockets.

![Client Test App](/assets/Screenshot_20250524_181321.png)

![Server Terminal](/assets/Screenshot_20250524_181419.png)


## ⚙️ Tech Stack

- **Runtime:** Node.js
- **Real-Time Communication:** Socket.IO
- **Server Framework:** Express.js
- **Environment Management:** dotenv

## 🚀 Setup Instructions

1. **Clone the Repository**

    ```bash
    git clone https://github.com/u-s1ddhar7h/internship-backend-tasks.git

    cd internship-backend-tasks/realtime-chat-backend
    ```

2. **Enter the Nix Dev Shell**

   ```bash
   nix develop
   ```

3. **Install Dependencies**

   ```bash
   npm install
   ```

4. **Configure Environment Variables**

    Create a `.env` file in the root directory:
    ```env
    PORT=4000
    JWT_SECRET=your_super_secret_32_characters_long_JWT_Token
    CLIENT_URL=http://localhost:3001
    ```

5. **Run the Server**

   ```bash
   npm run dev
   ```

    The server will run at:
    `ws://localhost:4000` (WebSocket)

## 🔌 WebSocket Events

The chat server communicates via custom Socket.IO events. Below are the key events supported:

### 🔄 Connection Events

| Event        | Description                         |
| ------------ | ----------------------------------- |
| `connect`    | Triggered when a client connects    |
| `disconnect` | Triggered when a client disconnects |

### 💬 Chat Events

| Event       | Payload                       | Description                          |
| ----------- | ----------------------------- | ------------------------------------ |
| `joinRoom`  | `{ roomId, username }`        | Joins a user to a specific chat room |
| `message`   | `{ roomId, message, sender }` | Sends a message to a room            |
| `leaveRoom` | `{ roomId, username }`        | Leaves the specified room            |
| `typing`    | `{ roomId, username }`        | Indicates that a user is typing      |

---
