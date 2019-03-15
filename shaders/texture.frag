#version 330
out vec4 frag_color;

in vec2 fragTexCoord;

uniform sampler2D tex;

void main() {
    frag_color = texture(tex, vec2(fragTexCoord.x, 1-fragTexCoord.y)); // Flip Y axis
}
