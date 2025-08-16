/** @type {import('stylelint').Config} */
module.exports = {
  extends: [
    "stylelint-config-standard",
  ],
  rules: {
    "property-no-unknown": [true, {
      ignoreProperties: ["/^--/"],
    }],

    "selector-pseudo-class-no-unknown": [true, {
      ignorePseudoClasses: ["global"],
    }],
  }
};
