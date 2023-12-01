const path = require('path');
const {glob} = require("glob");


module.exports = {
    entry: () => {
    	const entries = {};
    	const files = glob.sync('./src/**/*.+(js|ts)'); // Using glob to find js/ts files in src directory and subdirectories
    	files.forEach(file => {
      		const entryName = path.basename(file, path.extname(file)); // Extracting the entry name without the file extension
      		entries[entryName] = './' + file;
    	});
		console.log(entries);
    	return entries;
  	},
	module: {
		rules: [
      		{
        		test: /\.tsx?$/,
        		use: 'ts-loader',
        		exclude: /node_modules/,
      		},
    	],
  	},
	resolve: {
    	extensions: ['.tsx', '.ts', '.js'],
  	},
	output: {
    	filename: 'dest.js',
    	path: path.resolve(__dirname, 'public/js'),
  },
};