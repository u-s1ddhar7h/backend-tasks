require('dotenv').config();
const express = require('express');
const { createServer } = require('http');
const { Server } = require('socket.io');
const cors = require('cors');
const jwt = require('jsonwebtoken'); // Add this line

const app = express();
app.use(cors());

// Temporary route to get a test token - remove after testing
app.get('/test-token', (req, res) => {
    const token = jwt.sign(
        { id: 'test_user_id', username: 'tester' },
        process.env.JWT_SECRET,
        { expiresIn: '1h' }
    );
    res.send(token);
});

const httpServer = createServer(app);
const io = new Server(httpServer, {
  transports: ['websocket', 'polling'], // Explicitly declare transports
  cors: {
    origin: ["http://localhost:3001", "http://127.0.0.1:3001", "file://"], // Add all possible origins
    methods: ["GET", "POST"],
    credentials: true
  }
});

// Import socket handlers
require('./socket/handler.js')(io);

const PORT = process.env.PORT || 4000;
httpServer.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});