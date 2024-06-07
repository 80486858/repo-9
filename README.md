<div align="center">
  <a href="https://dash.plotly.com/project-maintenance">
    <img src="https://dash.plotly.com/assets/images/maintained-by-plotly.png" width="400px" alt="Maintained by Plotly">
  </a>
</div>

A simple fork of the [official Heroku Python buildpack](https://github.com/heroku/heroku-buildpack-python) to add a `postinstall` hook. This hook, activated by creating a shell script in the root of your repo called `postinstall`, will be run immediately after `pip install -r requirements.txt` but before compressing the slug. This allows you to uninstall packages that are hard dependencies of other packages you need, but aren't actually called when running your app.
