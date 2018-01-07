#version 330 core

in vec2 p_Texcoord;

out vec4 o_Color;

void main() {
    o_Color = vec4(p_Texcoord.x, 0.0, p_Texcoord.y, 1.0);
}
