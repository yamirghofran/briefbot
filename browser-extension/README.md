# BriefBot Browser Extension

A browser extension that allows you to submit URLs to your BriefBot reading list directly from any webpage.

## Features

- 🚀 **One-click URL submission** from any webpage
- 📝 **Auto-fill current page** URL and title
- 💾 **Persistent user ID** storage
- 🔔 **Success/error notifications**
- 🖱️ **Right-click context menu** support
- 📊 **Backend status monitoring**
- 🎨 **Beautiful, modern UI**

## Installation

### For Development (Unpacked Extension)

1. **Clone or download** this extension folder
2. **Open your browser's extension management page**:
   - **Chrome/Edge**: `chrome://extensions/` or `edge://extensions/`
   - **Firefox**: `about:debugging#/runtime/this-firefox`
3. **Enable Developer Mode** (Chrome/Edge only)
4. **Click "Load unpacked"** and select the extension folder
5. **Pin the extension** to your toolbar for easy access

### For Users (Production)

> ⚠️ **Note**: This extension is currently designed for development use. For production distribution, it would need to be published to the Chrome Web Store, Firefox Add-ons, or Edge Add-ons store.

## Usage

### Basic Usage

1. **Navigate to any webpage** you want to save
2. **Click the BriefBot extension icon** in your toolbar
3. **Enter your User ID** (get this from your BriefBot app)
4. **Click "Save"** to store your User ID
5. **Click "Submit to BriefBot"** to add the page to your reading list

### Context Menu (Right-click)

- **Right-click anywhere on a page** → Select "Submit to BriefBot"
- **Right-click on links** → Select "Submit to BriefBot" to submit the link URL

### Keyboard Shortcuts

- **Enter key** in User ID field: Saves User ID
- **Extension icon**: Opens popup for quick submission

## Configuration

### User ID Setup

Your User ID is stored locally in the extension. To find your User ID:

1. **Check your BriefBot application** - it's usually displayed in your user profile
2. **Look at the URL** when viewing your items - it often contains the user ID
3. **Check your browser's local storage** if you're a developer

### Backend Configuration

The extension connects to `http://localhost:8080` by default. To change this:

1. **Edit the `API_BASE_URL` in `popup.js`**
2. **Update the `host_permissions` in `manifest.json`**
3. **Reload the extension**

## File Structure

```
browser-extension/
├── manifest.json          # Extension configuration
├── popup.html            # Main popup interface
├── popup.css             # Styling for the popup
├── popup.js              # Main popup functionality
├── content.js            # Content script for page interaction
├── background.js         # Background service worker
├── icon-generator.js     # Icon creation utility
├── icons/                # Extension icons (to be created)
│   ├── icon16.png
│   ├── icon48.png
│   └── icon128.png
└── README.md             # This file
```

## Icons

The extension expects custom icons in the `icons/` folder:

- **icon16.png** (16x16 pixels) - Toolbar icon
- **icon48.png** (48x48 pixels) - Extension management page
- **icon128.png** (128x128 pixels) - Chrome Web Store

### Creating Icons

For development, you can create simple icons using any image editor:

1. **Create a square canvas** (16x16, 48x48, 128x128)
2. **Add a gradient background** (#667eea to #764ba2)
3. **Add white "BB" text** in a bold, centered font
4. **Save as PNG files**

Or use the provided `icon-generator.js` as a template.

## Browser Compatibility

- ✅ **Google Chrome** (Manifest V3)
- ✅ **Microsoft Edge** (Manifest V3)
- ✅ **Brave Browser** (Manifest V3)
- ⚠️ **Mozilla Firefox** (Manifest V2/V3 - may need adjustments)
- ⚠️ **Safari** (Requires additional configuration)

## Troubleshooting

### Common Issues

**Extension won't load:**
- Check that all required files are present
- Verify `manifest.json` syntax is valid
- Ensure file paths are correct

**Can't connect to backend:**
- Verify your BriefBot backend is running on `localhost:8080`
- Check CORS settings in your backend
- Ensure `host_permissions` in `manifest.json` includes your backend URL

**User ID not saving:**
- Check browser storage permissions
- Verify the User ID is a valid positive number
- Try reloading the extension

**No notifications appearing:**
- Check browser notification permissions
- Ensure `showNotifications` is enabled in storage

### Debug Mode

Open the browser's developer tools for the extension:

1. **Right-click the extension icon** → "Inspect popup"
2. **Check the Console tab** for error messages
3. **Check the Network tab** for API calls

## Security Notes

- The extension only communicates with `localhost:8080` by default
- User ID is stored locally in browser storage
- No sensitive data is transmitted without user interaction
- All API calls use the same authentication as your main BriefBot app

## Development

### Making Changes

1. **Edit the relevant files**
2. **Reload the extension** in browser extension management
3. **Test your changes**
4. **Check browser console for errors**

### Adding Features

The extension is modular:
- **Popup UI**: Edit `popup.html`, `popup.css`, `popup.js`
- **Background tasks**: Edit `background.js`
- **Page interaction**: Edit `content.js`
- **Configuration**: Edit `manifest.json`

## License

This extension is designed to work with BriefBot. Check your BriefBot application's license for usage terms.

## Support

For issues specific to this extension:
1. Check the troubleshooting section above
2. Review browser console logs
3. Verify your BriefBot backend is running correctly

For BriefBot application issues, consult your BriefBot documentation.