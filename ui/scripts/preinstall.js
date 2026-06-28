require('./loadPlugins'); 
const execPath = process.env.npm_execpath || '';

if ( !/npm|pnpm|yarn|bun/.test(execPath) ) {
    console.warn( '\u001b[33mUnsupported package manager.\u001b[39m\n', ); 
    process.exit(1);
}