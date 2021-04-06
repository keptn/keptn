const express = require('express');
const router = express.Router();

/**
 * Router level middleware for login
 */
router.get('/login', (req, res, next) => {
  console.log('Login to Keptn bridge.');

  // todo replace with redirects to OAUTH SERVICE
  req.session.authenticated = true;

  res.redirect('/');
  return res;
});

/**
 * Router level middleware for logout
 */
router.get('/logout', (req, res) => {
  console.log('Logout from Keptn bridge.');

  req.session.destroy(function (err) {
    res.redirect('/');
    return res;
  });

});

module.exports.router = router;
