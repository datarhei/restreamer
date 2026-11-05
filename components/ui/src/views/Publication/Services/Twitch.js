import React from 'react';

import { faTwitch } from '@fortawesome/free-brands-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Trans } from '@lingui/macro';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';

import FormInlineButton from '../../../misc/FormInlineButton';
import Select from '../../../misc/Select';

const id = 'twitch';
const name = 'Twitch';
const version = '1.1';
const stream_key_link = 'https://dashboard/settings/stream';
const description = <Trans>Live-Streaming to Twitch Live RTMP Service.</Trans>;
const image_copyright = <Trans>More about licenses here</Trans>;
const author = {
	creator: {
		name: 'datarhei',
		link: 'https://github.com/datarhei',
	},
	maintainer: {
		name: 'datarhei',
		link: 'https://github.com/datarhei',
	},
};
const category = 'platform';
const requires = {
	protocols: ['rtmp'],
	formats: ['flv'],
	codecs: {
		audio: ['aac', 'mp3'],
		video: ['h264'],
	},
};

function ServiceIcon(props) {
	return <FontAwesomeIcon icon={faTwitch} style={{ color: '#9147FF' }} {...props} />;
}

function init(settings) {
	const initSettings = {
		region: 'live',
		key: '',
		...settings,
	};

	return initSettings;
}

function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		settings[what] = value;

		const output = createOutput(settings);

		props.onChange([output], settings);
	};

	const createOutput = (settings) => {
		let region_postfix = '.twitch.tv';
		if (settings.region.includes('live-video.net')) {
			region_postfix = '';
		}

		const output = {
			address: 'rtmp://' + settings.region + region_postfix + '/app/' + settings.key,
			options: ['-f', 'flv'],
		};

		return output;
	};

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<Select label={<Trans>Region</Trans>} value={settings.region} onChange={handleChange('region')}>
					<MenuItem value="live">Autodetect</MenuItem>
					<MenuItem value="live-hkg">Asia - Hong Kong</MenuItem>
					<MenuItem value="blr01.contribute.live-video.net">Asia - India, Bangalore</MenuItem>
					<MenuItem value="maa01.contribute.live-video.net">Asia - India, Chennai</MenuItem>
					<MenuItem value="hyd01.contribute.live-video.net">Asia - India, Hyderabad</MenuItem>
					<MenuItem value="bom01.contribute.live-video.net">Asia - India, Mumbai</MenuItem>
					<MenuItem value="del01.contribute.live-video.net">Asia - India, New Delhi</MenuItem>
					<MenuItem value="jkt01.contribute.live-video.net">Asia - Indonesia, Cikarang Barat</MenuItem>
					<MenuItem value="jkt02.contribute.live-video.net">Asia - Indonesia, Jakarta</MenuItem>
					<MenuItem value="osa01.contribute.live-video.net">Asia - Japan, Osaka</MenuItem>
					<MenuItem value="live-tyo">Asia - Japan, Tokyo</MenuItem>
					<MenuItem value="mnl01.contribute.live-video.net">Asia - Philippines, Manila</MenuItem>
					<MenuItem value="live-sin">Asia - Singapore</MenuItem>
					<MenuItem value="live-sel">Asia - South Korea, Seoul</MenuItem>
					<MenuItem value="live-tpe">Asia - Taiwan, Taipei</MenuItem>
					<MenuItem value="bkk02.contribute.live-video.net">Asia - Thailand, Bangkok</MenuItem>
					<MenuItem value="live-vie">Europe - Austria, Vienna</MenuItem>
					<MenuItem value="live-prg">Europe - Czech Republic, Prague</MenuItem>
					<MenuItem value="live-cph">Europe - Denmark, Copenhagen</MenuItem>
					<MenuItem value="live-hel">Europe - Finland, Helsinki</MenuItem>
					<MenuItem value="live-mrs">Europe - France, Marseille</MenuItem>
					<MenuItem value="live-cdg">Europe - France, Paris</MenuItem>
					<MenuItem value="live-ber">Europe - Germany, Berlin</MenuItem>
					<MenuItem value="dus01.contribute.live-video.net">Europe - Germany, Dusseldorf</MenuItem>
					<MenuItem value="live-fra">Europe - Germany, Frankfurt</MenuItem>
					<MenuItem value="muc01.contribute.live-video.net">Europe - Germany, Munich</MenuItem>
					<MenuItem value="live-mil">Europe - Italy, Milan</MenuItem>
					<MenuItem value="live-ams">Europe - Netherlands, Amsterdam</MenuItem>
					<MenuItem value="live-osl">Europe - Norway, Oslo</MenuItem>
					<MenuItem value="live-waw">Europe - Poland, Warsaw</MenuItem>
					<MenuItem value="live-lis">Europe - Portugal, Lisbon</MenuItem>
					<MenuItem value="live-mad">Europe - Spain, Madrid</MenuItem>
					<MenuItem value="live-arn">Europe - Sweden, Stockholm</MenuItem>
					<MenuItem value="live-lhr">Europe - UK, London</MenuItem>
					<MenuItem value="live-ymq">NA - Canada, Quebec</MenuItem>
					<MenuItem value="live-yto">NA - Canada, Toronto</MenuItem>
					<MenuItem value="live-qro">NA - Mexico, Queretaro</MenuItem>
					<MenuItem value="live-syd">Oceania - Australia, Sydney</MenuItem>
					<MenuItem value="live-scl">South America - chile, Santiago</MenuItem>
					<MenuItem value="for01.contribute.live-video.net">South America - Brazil, Fortaleza</MenuItem>
					<MenuItem value="live-rio">South America - Brazil, Rio de Janeiro</MenuItem>
					<MenuItem value="live-sao">South America - Brazil, Sao Paulo</MenuItem>
					<MenuItem value="live-eze">South America - Buenos Aires, Argentina</MenuItem>
					<MenuItem value="bog01.contribute.live-video.net">South America - Colombia, Bogota</MenuItem>
					<MenuItem value="live-mde">South America - Colombia, Medellin</MenuItem>
					<MenuItem value="live-lim">South America - Peru, Lima</MenuItem>
					<MenuItem value="live-dfw">US Central - Dallas, TX</MenuItem>
					<MenuItem value="live-den">US Central - Denver, CO</MenuItem>
					<MenuItem value="dfw56.contribute.live-video.net">US Central - Garland, TX</MenuItem>
					<MenuItem value="live-hou">US Central - Houston, TX</MenuItem>
					<MenuItem value="live-slc">US Central - Salt Lake City, UT</MenuItem>
					<MenuItem value="live-iad">US East - Ashburn, VA</MenuItem>
					<MenuItem value="live-atl">US East - Atlanta, GA</MenuItem>
					<MenuItem value="live-ord">US East - Chicago, IL</MenuItem>
					<MenuItem value="mfe01.contribute.live-video.net">US East - McAllen, TX</MenuItem>
					<MenuItem value="live-mia">US East - Miami, FL</MenuItem>
					<MenuItem value="live-jfk">US East - New York, NY</MenuItem>
					<MenuItem value="live-lax">US West - Los Angeles, CA</MenuItem>
					<MenuItem value="live-phx">US West - Phoenix, AZ</MenuItem>
					<MenuItem value="live-pdx">US West - Portland, Oregon</MenuItem>
					<MenuItem value="slc.contribute.live-video.net">US West - Salt Lake City, UT</MenuItem>
					<MenuItem value="live-sfo">US West - San Francisco, CA</MenuItem>
					<MenuItem value="live-sjc">US West - San Jose, CA</MenuItem>
					<MenuItem value="live-sea">US West - Seattle, WA</MenuItem>
				</Select>
			</Grid>
			<Grid item xs={12} md={9}>
				<TextField variant="outlined" fullWidth label={<Trans>Stream key</Trans>} value={settings.key} onChange={handleChange('key')} />
			</Grid>
			<Grid item xs={12} md={3}>
				<FormInlineButton target="blank" href={stream_key_link} component="a">
					<Trans>GET</Trans>
				</FormInlineButton>
			</Grid>
		</Grid>
	);
}

Service.defaultProps = {
	settings: {},
	skills: {},
	metadata: {},
	streams: [],
	onChange: function (output, settings) {},
};

export { id, name, version, stream_key_link, description, image_copyright, author, category, requires, ServiceIcon as icon, Service as component };
