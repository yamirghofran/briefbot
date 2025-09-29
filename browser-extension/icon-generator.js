// Simple icon generator for BriefBot extension
// This creates basic colored squares with text for development

// Create a simple canvas-based icon
function createIcon(text, size, color) {
    const canvas = document.createElement('canvas');
    canvas.width = size;
    canvas.height = size;
    const ctx = canvas.getContext('2d');
    
    // Background
    ctx.fillStyle = color;
    ctx.fillRect(0, 0, size, size);
    
    // Text
    ctx.fillStyle = 'white';
    ctx.font = `bold ${size * 0.6}px Arial`;
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(text, size / 2, size / 2);
    
    return canvas.toDataURL();
}

// For now, let's create placeholder instructions
console.log(`
BriefBot Extension Icons

Since this is a development environment, here are simple instructions to create icons:

1. Create 3 square PNG images (16x16, 48x48, 128x128 pixels)
2. Use a blue/purple gradient background (#667eea to #764ba2)
3. Add white "BB" text in the center
4. Save them as:
   - icon16.png (16x16)
   - icon48.png (48x48) 
   - icon128.png (128x128)

Alternative: Use any image editor to create simple icons with:
- Background: Linear gradient from #667eea to #764ba2
- Text: "BB" in white, bold font
- Centered text

The extension will work without custom icons (browser will use default puzzle piece),
but custom icons make it easier to identify in the toolbar.
`);