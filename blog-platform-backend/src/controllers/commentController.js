const Comment = require('../models/Comment');

const createComment = async (req, res) => {
  try {
    const { content } = req.body;
    const commentId = await Comment.create({
      content,
      userId: req.user.id,
      blogId: req.params.blogId
    });
    
    const comments = await Comment.findByBlogId(req.params.blogId);
    res.status(201).json(comments);
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

const getBlogComments = async (req, res) => {
  try {
    const comments = await Comment.findByBlogId(req.params.blogId);
    res.json(comments);
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

const deleteComment = async (req, res) => {
  try {
    // In a real app, you might want to check if the user owns the comment
    await Comment.delete(req.params.commentId);
    res.json({ message: 'Comment deleted successfully' });
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

module.exports = {
  createComment,
  getBlogComments,
  deleteComment
};