# github.com/magiccap/magiccap-uploaders-kernel
The kernel for MagicCap uploaders.

## Adding uploaders to the MagicCap kernel
We welcome the following into the kernel:
- Open standards (such as POMF and S3).
- Support for uploaders with a reasonably sized user base behind them.

To add your uploader, simply check if the standard it is based on (for example HTTP) is already added. You can check this in the `standards` folder. If/when it is, you can go into the `uploaders` folder and create a folder for the standard. From here, insert the JSON into the folder. You can use other uploaders as an example, just make sure that the `spec` object is specific to the standard you are using (this is standard dependant). From here, you can edit `v1_imports.json`.

If you do not mind tests to your uploader being automatically ran, you can make a test in the `tests` folder. Note that if you need a secret added, you will need an administrator in the repository to add it. If you need it, just ask!

When a uploader is added/updated and pushed into master, it will be added into a JSON file which will be downloaded by MagicCap 3.0+ clients roughly every 10 minutes. This means that you do not need to worry about the user needing to update MagicCap unless you are adding a new standard.

## Building for importing outside of the official repository
If you want to do this, simply copy the `uploaders` folder. From here, you can delete all the implementations/icons you don't need (**DON'T FORGET TO REMOVE THEM FROM `v1_imports.json` TOO!**) and just add what you need. From here, running `build.py` will give you a file you can import.
