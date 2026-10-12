#version 300 es

precision mediump float;

uniform vec2 uResolution;
uniform float uTime;

out vec4 outColor;


float random(in vec2 st){

	return fract(sin(dot(st.xy, vec2(12.9898,78.233))) * 43758.5452122);
}

float snoise(in vec2 st){

	vec2 i = floor(st);
	vec2 f = fract(st);
	
	float a = random(i);
	float b = random(i+vec2(1.0,0.0));
	float c = random(i+vec2(0.0,1.0));
	float d = random(i+vec2(1.0,1.0));
	
	vec2 u = f * f * (3.0-2.0 * f);
	
	return((mix(a, b, u.x) +
			(c-a) * u.y * (1.0-u.x) +
			(d-b) * u.x * u.y)*2.0)-1.0;
}

float noise(in vec2 st){

	vec2 i = floor(st);
	vec2 f = fract(st);
	
	float a = random(i);
	float b = random(i+vec2(1.0,0.0));
	float c = random(i+vec2(0.0,1.0));
	float d = random(i+vec2(1.0,1.0));
	
	vec2 u = f * f * (3.0-2.0 * f);
	
	return mix(a, b, u.x) +
			(c-a) * u.y * (1.0-u.x) +
			(d-b) * u.x * u.y;
}

#define OCTAVES 7
float fbm(in vec2 st){

	//float value = 0.0;
	//float amplitude =.5;
	//float frequency = 0.;
	//float lacunarity = 2.1;
	//for (int i = 0; i< OCTAVES; i++){
	//	value += amplitude * abs(snoise(st));
	//	st *= lacunarity;
	//	amplitude*= .5;
	//
	//}
	//return value;
	
	float value = 0.0;
	float amplitude =.5;
	float frequency = 0.;
	float lacunarity = 2.1;
	vec2 shift = vec2(100.0,100.0);
	
	// rotate to reduce axial bias
	mat2 rot = mat2(cos(0.5), sin(0.5), 
					-sin(0.5), cos(0.50));
	
	for (int i = 0; i< OCTAVES; i++){
		value += amplitude * abs(snoise(st));
		st = st * rot * 2.0 + shift;
		amplitude*= .5;
	
	}
	return value;
}

float warpedfbm(in vec2 st, in float offset){

   	vec2 q = vec2(0.0,0.0);
    q.x = fbm(st + 0.01*offset);
    q.y = fbm(st + vec2(1.0,1.0));

	vec2 r = vec2(0.0);
	r.x = fbm(st + 1.0*q + vec2(1.7,9.2)+ 0.15*offset);
    r.y = fbm(st + 1.0*q + vec2(8.3,2.8)+0.126*offset);

	return fbm(st+r);
}

void main()
{
    vec2 uv = gl_FragCoord.xy/uResolution.xy;
    uv.x *= uResolution.x/uResolution.y;
    uv.y *= 1.5;
    
    	
   vec2 delta = vec2(uTime,uTime)*0.0;
    
   float r = warpedfbm((uv*2.0), uTime*0.1);
   float g = warpedfbm((uv*2.5), uTime*0.6);
   float b = warpedfbm((uv*0.8), uTime*0.2);
    
    
    outColor = vec4(r,g,b,1.0);
}