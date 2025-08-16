import esbuild from 'esbuild';
import fs from 'fs';

const appEntries = {
  app: 'web/app/ts/app.ts',
  login: 'web/login/ts/login.ts',
};

const cssEntries = {
  app: 'web/app/css/app.css',
  login: 'web/login/css/login.css',
};

const swDir = 'web/sw';
const swEntries = fs.readdirSync(swDir)
  .filter(f => f.endsWith('.sw.ts'))
  .reduce((acc, f) => {
    const name = f.replace(/\.ts$/, ''); // e.g., auth.sw
    acc[name] = `${swDir}/${f}`;
    return acc;
  }, {});

async function build() {
  try {
    await esbuild.build({
      entryPoints: appEntries,
      bundle: true,
      minify: true,
      sourcemap: true,
      outdir: 'static/js',
      entryNames: '[name]',
      format: 'esm',

    });

    await esbuild.build({
      entryPoints: cssEntries,
      bundle: true,
      minify: true,
      sourcemap: true,
      outdir: 'static/css',
      entryNames: '[name]',
      resolveExtensions: ['.css'],
    });

    if (Object.keys(swEntries).length > 0) {
      await esbuild.build({
        entryPoints: swEntries,
        bundle: true,
        minify: true,
        sourcemap: true,
        outdir: 'static/js/sw',
        entryNames: '[name]',
        format: 'esm',
      });
    }

    console.log('JS, SW, and CSS builds succeeded.');
  } catch (err) {
    console.error('Build failed:', err);
    process.exit(1);
  }
}

build();
