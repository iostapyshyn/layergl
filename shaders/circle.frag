#version 330
out vec4 frag_color;

uniform vec4 color;
uniform vec3 circle;

const int aa = 1;

void main() {
    float d = distance(gl_FragCoord.xy, circle.xy);
    float k = 1-smoothstep(circle.z-(2*aa), circle[2], d);
    
    frag_color = vec4(color.xyz, k*color.w);
}
