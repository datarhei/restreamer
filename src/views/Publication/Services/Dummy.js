import React from 'react';
import PropTypes from 'prop-types';

import { Trans } from '@lingui/macro';
import ExtensionIcon from '@mui/icons-material/Extension';
import Grid from '@mui/material/Grid';
import TextField from '@mui/material/TextField';

import Checkbox from '../../../misc/Checkbox';

// id is a unique identifier for this service.
const id = 'dummy';

// name is the name of the service as it will be displayed in the UI.
const name = 'Dummy Live';

// version is the version of this service and it will be displayed
// in the UI.
const version = '1.0';

// link to the stream key
const stream_key_link = '';

// description is a brief description of this service and it will be
// displayed in the UI.
const description = <Trans>This is a dummy service that explains to you the concepts of service.</Trans>;

// image_copyright is a short explanation regarding the copyrights for
// the target of this service. Ideally it should contain also a link
// to a webpage with more information.
const image_copyright = <Trans>This is to mention the copyright regulations for the target of this service.</Trans>;

// author mentiones the author of this service.
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

// category is the category this service belongs to. It has to be
// one of 'platform', 'universal', 'software'. Choose 'platform' if
// the target of this service is a distribution platform. Choose
// 'universal' if this service just distributes the data via a certain
// protocol without targeting a specific platform. Choose 'software'
// if this service targets a specific streaming software.
const category = 'platform';

// requires is an object holding the requirements for this service.
// Depending on what is available from ffmpeg, the service may not be
// listed. The full form of this object is:
//
//    {
//       protocols: array of required protocols
//       formats: array of required formats
//       codecs: {
//	        audio: array of required audio codecs
//          video: array of required video codecs
//       }
//    }
//
// The 'protocols' array holds a list of output protocols of which this
// service needs at least one to be available in order to operate
// properly. If the array is empty or the 'protocols' field is missing,
// then there are no requirements regarding the available protocols.
// The available protocols are derived from the -protocols output of
// ffmpeg.
//
// Example: ['rtmp', 'rtmps']
//
// The 'formats' array holds a list of output formats of which this service
// needs at least one to be available in order to operate properly. If
// the array is empty or the 'formats' field is missing, then there are no
// requirements regarding the available formats. The available formats
// are derived from the -formats output of ffmpeg.
//
// Example: ['flv']
//
// The 'codecs' objects holds the fields 'audio' and 'video' which are
// arrays that hold a list of encoding codecs. At least one of the codecs
// needs to be available for this service to be able to operate properly.
// If the arrays are empty or if the 'audio' and/or 'video' fields are
// missing, then there are no requirements regarding the available codecs.
// If neither audio nor video codecs are required, the 'codecs' field
// can be omitted. The available codecs are derived from the -codecs
// output of ffmpeg.
//
// Example: { audio: ['aac'], video: ['h264'] }
//
// Usually only 'protocols' and 'formats' have some requirements.
// Only if you need to do some re-encoding, you have to specifiy the
// required codecs.
//
// If one of the requirements cannot be fulfilled, then this service
// will not be listed and it will not be possible to add a new publication
// using this service and this service will not be rendered. However, if
// there's already a publication using this service that has been created
// with different requirements or based on a different ffmpeg binary, then
// this service might still be called. In this case, the prop 'skills'
// will be 'null' and you can inform the user accordingly.

// The requirements are not fulfilled if one of the
// requirements doesn't have at least one match. For example, you have
// the requirement:
//
//    {
//       protocols: ['rtmp', 'rtmps'],
//       formats: ['flv'],
//    }
//
// If neither the protocol 'rtmp' nor 'rtmps' is available or the format
// 'flv' is not available, then the requirements are not fulfilled. In
// contrast, if only the protocol 'rtmps' is not available, then requirements
// are fulfilled, because at least the protocol 'rtmp' and the format 'flv'
// are available.
//
// Based on the specified requirements the 'skills' prop will contain which
// of the required protocols, formats, and/or codecs are available. The
// 'skills' prop are described below.
//
// Note that you only specify the wanted codecs (i.e.
// for H264 encoding you list 'h264' and not 'libx264' )
const requires = {
	protocols: ['dummy'],
	formats: ['flv'],
};

// ServiceIcon is a React component and returns the icon for this service
function ServiceIcon(props) {
	return <ExtensionIcon style={{ color: '#FFFFFF' }} {...props} />;
}

// Initialize the settings, i.e. set defaults for fields that are not defined
function init(settings) {
	const initSettings = {
		stream_key: '',
		rtmp_primary: true,
		rtmp_backup: false,
		...settings,
	};

	return initSettings;
}

// createOutput creates the FFmpeg output options based on the settings,
// skills, metadata, and/or streams. It returns an array of arrays of options
// for each output.
//
// This function should be exported, such that it can be called from the
// outside in case e.g. the streams changed (due to changes in the encoding
// settings) and the output options depend on that information.
function createOutputs(settings, skills, metadata, streams) {
	settings = init(settings);
	const outputs = [];

	if (settings.stream_key.length === 0) {
		return outputs;
	}

	if (settings.rtmp_primary) {
		outputs.push({
			address: 'rtmp://a.rtmp.youtube.com/live2/' + settings.stream_key,
			options: ['-codec', 'copy', '-f', 'flv'],
		});
	}

	if (settings.rtmp_backup) {
		outputs.push({
			address: 'rtmp://b.rtmp.youtube.com/live2?backup=1/' + settings.stream_key,
			options: ['-codec', 'copy', '-f', 'flv'],
		});
	}

	return outputs;
}

// Service is a React component that implements a service.
//
// A service receives these props:
//
//   - settings: object
//   - skills: object
//   - metadata: object
//   - onChange: function
//
// The contents of the 'settings' object are completely up to the
// service. It holds all settings for this service. The state of the
// settings is managed by the parent React component.
//
// The 'skills' object has information about the available protocols,
// formats, and codecs based on your requirements. It has to form:
//
//   {
//      protocols: array of available protocols
//      formats: array of available formats
//      codecs: {
//         audio: {
//            [codec]: array of available coders,
//            ...
//         },
//         video: {
//            [codec]: array of available coders,
//            ...
//         },
//      }
//   }
//
// The 'protocols' array lists all the protocols from your requirments
// that are actually available. The array is empty if you don't have any
// protocol requirements.
//
// The 'formats' array lists all the formats from your requirements that
// are actually available. The array is empty if you don't have any format
// requirements.
//
// The 'codecs' object holds an 'audio' and a 'video' object. These objects
// have the available codecs from the requirements as fields. Each codec
// has a list of available coders. For example for the codec 'h264' there
// are the coders 'libx264' and 'h264_nvenc' available where the service
// or user has to choose from. This depends on with which libraries ffmpeg
// has been compiled with.
//
// The 'metadata' object holds metadata about the stream that this
// service will process. It has this layout:
//
//   {
//      name: string,
//      description: string,
//      license: string,
//   }
//
// The 'onChange' has the signature
//
//   function(outputs: array, settings: object)
//
// It has to be called whenever there are changes (i.e. user input) of
// the settings. 'outputs' is an array of output definitions for ffmpeg
// of the form:
//
//   {
//      address: string,
//      options: array,
//   }
//
// Each output requires an address and a list of options. These options
// are the ffmpeg command line options that are necessary for this service
// and are based on the settings. The parent React component is not
// aware of the internals of this service.
//
// 'settings' is the settings object where you store the settings for the
// service. Its state is managed by the parent React component.
function Service(props) {
	const settings = init(props.settings);

	const handleChange = (what) => (event) => {
		const value = event.target.value;

		if (['rtmp_primary', 'rtmp_backup'].includes(what)) {
			settings[what] = !settings[what];
		} else {
			settings[what] = value;
		}

		const outputs = createOutputs(settings, props.skills, props.metadata, props.streams);

		props.onChange(outputs, settings);
	};

	if (props.skills === null) {
		return null;
	}

	return (
		<Grid container spacing={2}>
			<Grid item xs={12}>
				<TextField variant="outlined" fullWidth label={<Trans>Stream key</Trans>} value={settings.stream_key} onChange={handleChange('stream_key')} />
			</Grid>
			<Grid item xs={12}>
				<Checkbox label={<Trans>Enable primary stream</Trans>} checked={settings.rtmp_primary} onChange={handleChange('rtmp_primary')} />
				<Checkbox label={<Trans>Enable backup stream</Trans>} checked={settings.rtmp_backup} onChange={handleChange('rtmp_backup')} />
			</Grid>
		</Grid>
	);
}

Service.defaultProps = {
	settings: {},
	skills: null,
	metadata: {},
	streams: [],
	onChange: function (output, settings) {},
};

Service.propTypes = {
	metadata: PropTypes.object.isRequired,
	onChange: PropTypes.func.isRequired,
	settings: PropTypes.object.isRequired,
	skills: PropTypes.object,
	streams: PropTypes.array.isRequired,
};

export {
	id,
	name,
	version,
	stream_key_link,
	description,
	image_copyright,
	author,
	category,
	requires,
	ServiceIcon as icon,
	Service as component,
	createOutputs,
};
