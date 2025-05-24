const express = require('express');
const router = express.Router();
const commentController = require('../controllers/commentController');
const auth = require('../middleware/auth');

router.post('/:blogId', auth, commentController.createComment);
router.get('/:blogId', commentController.getBlogComments);
router.delete('/:commentId', auth, commentController.deleteComment);

module.exports = router;