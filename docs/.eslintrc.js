const importRules = require('./.eslintrc-rules-import')

module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  parser: 'vue-eslint-parser',
  parserOptions: {
    parser: '@typescript-eslint/parser',
    ecmaVersion: 'latest',
    sourceType: 'module',
  },
  plugins: ['@typescript-eslint', 'import'],
  extends: ['eslint:recommended', 'plugin:import/recommended', 'prettier'],
  overrides: [
    {
      files: ['.vuepress/**/*.ts'],
      extends: [
        'eslint:recommended',
        'plugin:@typescript-eslint/recommended',
        'plugin:import/recommended',
        'plugin:import/typescript',
        'prettier',
      ],
    },
    {
      files: ['.vuepress/**/*.vue'],
      extends: [
        'eslint:recommended',
        'plugin:@typescript-eslint/recommended',
        'plugin:vue/vue3-recommended',
        'plugin:import/recommended',
        'plugin:import/typescript',
        'prettier',
      ],
    },
  ],
  settings: {
    'import/parsers': {
      '@typescript-eslint/parser': ['.ts', '.tsx'],
      'vue-eslint-parser': ['.vue'],
    },
    'import/resolver': {
      typescript: {
        alwaysTryTypes: true,
      },
    },
  },
  rules: {
    'no-debugger': 'off',
    curly: ['error', 'all'],
    'prefer-const': 'error',
    'no-var': 'error',
    '@typescript-eslint/no-empty-function': 'off',
    ...importRules,
  },
}
