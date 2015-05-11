html
<p>It was three full days into a five day training course when suddenly the students - Support staff - became very animated and excited.</p>

<p>They were watching all the aggregated Cloud Foundry component logs &amp; events from across the entire system.</p>

<p>"Oh I wish you'd shown us this on day one!"</p>

<p>Ahh blessed Support people can be more excited about what's wrong with a system than the incredible discovery that something like Cloud Foundry even works at all.</p>

<p>Cloud Foundry includes two layers of logging:</p>

<ul>
<li>user/application logs and events</li>
<li>operations component logs and message bus events</li>
</ul>

<p>The former are expressly for the end users/teams. These logs can be streamed to the local command line or continously drained off to logging services like Papertrail or Logentries.</p>

<p>The latter was what the Support training students had seen. All the components' logs (router, cloudcontroller, health manager, runner, etc) and all the intercommunication messages (message bus) were being aggregated and stored.</p>

<p>In the demo to the students, and the video below, we use the built-in example aggregator that comes bundled with Cloud Foundry. In production, most operations people connect in Logstash or Splunk or similar.</p>

<div class='embed-container'>  
<iframe  src='//www.youtube.com/embed/WCvxSBigmog' frameborder='0' allowfullscreen></iframe>  
</div>

<h2 id="enablingsyslogforcomponents">Enabling syslog for components</h2>

<p>Whether you use the built-in example aggregator job or a remote syslog receiver such as Logstash or Splunk, you want to add the following global property to your BOSH deployment manifest for Cloud Foundry.</p>

<pre><code class="yml">properties:  
  syslog_aggregator:
    address: 10.15.213.19
    port: 54321
</code></pre>

<p>Where the <code>address</code> and <code>port</code> values are for the inbound endpoint of your logging aggregator.</p>

<p>By default <code>tcp</code> transport is used. You can change it to <code>udp</code> or <code>repl</code>:</p>

<pre><code class="yml">properties:  
  syslog_aggregator:
    address: 10.15.213.19
    port: 54321
    transport: udp
</code></pre>

<h2 id="runningtheexampleaggregator">Running the example aggregator</h2>

<p>In the video above we run the <code>syslog_aggregator</code> job that is built-in with Cloud Foundry's BOSH release.</p>

<p>Example configuration for running it could be:</p>

<pre><code class="yml">jobs:  
  - name: syslog_aggregator_z1
    templates:
      - name: syslog_aggregator
        release: cf
    instances: 1
    resource_pool: small_z1
    networks:
      - name: cf1
        default: [dns, gateway]
        static_ips:
          - 10.15.213.19
</code></pre>

<p>The assigned <code>static_ip</code> should be used by the <code>properties.syslog_aggregator.address</code> configuration in the previous section.</p>

<p>CODE: LETTUCE</p>
