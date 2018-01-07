#version 330 core

uniform mat4 u_Model;
uniform mat4 u_View;
uniform mat4 u_Projection;

in vec3 a_Vertex;
in vec2 a_Texcoord;

out vec2 p_Texcoord;

void main() {
    p_Texcoord = a_Texcoord;
    gl_Position = u_Projection * u_View * u_Model * vec4(a_Vertex, 1.0);
}
