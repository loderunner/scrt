/** @type {import('vls').VeturConfig} */
module.exports = {
  settings: {
    'vetur.useWorkspaceDependencies': true,
  },
  projects: [
    './docs', // Shorthand for specifying only the project root location
    {
      root: './docs',
    },
  ],
};
